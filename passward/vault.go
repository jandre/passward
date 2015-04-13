package passward

import (
	"crypto/rand"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/jandre/passward/util"
	"github.com/jandre/sshcrypt"
)

type EncryptedValue struct {
	Path string
}

type Entry struct {
	Name   string
	Path   string
	Values map[string]EncryptedValue
}

type Vault struct {
	Name           string
	Description    string
	RemoteUpstream string
	Path           string           `toml:"-"`
	Entries        map[string]Entry `toml:"-"`
	users          *VaultUsers      `toml:"-"`
	git            *Git             `toml:"-"`
	username       string           `toml:"-"`
	email          string           `toml:"-"`
}

type VaultUsers struct {
	Path  string
	users map[string]*VaultUser
}

type VaultUser struct {
	path            string
	email           string
	publicKeyString string
	encryptedKey    string
	publicKey       sshcrypt.PublicKey
}

func (vu *VaultUser) Save() error {
	if !util.DirectoryExists(vu.path) {
		if err := os.MkdirAll(vu.path, 0700); err != nil {
			return err
		}
	}
	keyfile := vu.publicKeyFile()
	if err := ioutil.WriteFile(keyfile, []byte(vu.publicKeyString), 0600); err != nil {
		return err
	}

	encryptedMaster := vu.encryptedMasterFile()
	if err := ioutil.WriteFile(encryptedMaster, []byte(vu.encryptedKey), 0600); err != nil {
		return err
	}
	return nil
}

func (vu *VaultUser) encryptedMasterFile() string {
	return path.Join(vu.path, "encrypted_master")
}

func (vu *VaultUser) publicKeyFile() string {
	return path.Join(vu.path, "key")
}

func NewVaultUser(usersPath string, email string, publicKey string, encryptedKey string) (*VaultUser, error) {
	var user VaultUser
	var err error
	user.path = path.Join(usersPath, email)
	user.email = email
	user.encryptedKey = encryptedKey
	user.publicKeyString = publicKey

	user.publicKey, err = sshcrypt.ParsePublicKey([]byte(publicKey))

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func ReadVaultUser(pathToUser string) (*VaultUser, error) {
	// TODO email := path.Base(pathToUser)
	return nil, nil
}

func NewVaultUsers(parentPath string) *VaultUsers {
	vaultUsersFolder := path.Join(parentPath, "users")
	result := VaultUsers{Path: vaultUsersFolder, users: make(map[string]*VaultUser, 0)}
	return &result
}

func (vu *VaultUsers) AddUser(email string,
	publicKeyString string,
	encryptedKey string) error {

	return nil
}

func (vu *VaultUsers) Initialize() error {
	if !util.DirectoryExists(vu.Path) {
		os.MkdirAll(vu.Path, 0700)
		ioutil.WriteFile(path.Join(vu.Path, ".placeholder"), nil, 0700)
	}
	return nil
	// TODO: read users
}

func ReadAllVaults(vaultPath string) (map[string]*Vault, error) {
	vaults := make(map[string]*Vault, 0)

	files, err := ioutil.ReadDir(vaultPath)

	if err != nil {
		return nil, err
	}

	for _, name := range files {
		vault, err := ReadVault(vaultPath, name.Name())
		if err != nil {
			return nil, err
		}

		if vault != nil {
			vaults[name.Name()] = vault
		}
	}

	return vaults, nil
}

func ReadVault(vaultPath string, name string) (*Vault, error) {
	dst := path.Join(vaultPath, name)
	configPath := filepath.Join(dst, "config.toml")

	var vault Vault
	_, err := toml.DecodeFile(configPath, &vault)

	if err != nil {
		return nil, err
	}

	vault.Path = dst // in case it was moved
	vault.users = NewVaultUsers(dst)
	vault.git = NewGit(dst, vault.username, vault.email)
	vault.Initialize()
	return &vault, nil
}

//
// Create a new vault.
//
func NewVault(vaultPath string, name string, username string, email string) (*Vault, error) {
	dst := path.Join(vaultPath, name)

	result := Vault{Name: name,
		Path:     dst,
		users:    NewVaultUsers(dst),
		email:    email,
		username: username,
		git:      NewGit(dst, username, email),
	}

	if err := result.Initialize(); err != nil {
		return nil, err
	}
	return &result, nil
}

func (v *Vault) configPath() string {
	return path.Join(v.Path, "config.toml")
}

func (v *Vault) saveConfig() error {
	file, err := os.Create(v.configPath())
	if err != nil {
		return err
	}
	defer file.Close()
	if err := toml.NewEncoder(file).Encode(v); err != nil {
		return err
	}
	return nil

}

//
// Initialize the new vault by performing a git init at path, etc.
//
func (v *Vault) Initialize() error {
	var err error

	// it's already setup
	if util.DirectoryExists(v.Path) {
		if err = v.git.Initialize(); err != nil {
			return err
		}
		if err = v.users.Initialize(); err != nil {
			return err
		}
	} else {

		debug("initializing vault: %s", v.Path)

		if err = v.git.Initialize(); err != nil {
			return err
		}
		if err = v.users.Initialize(); err != nil {
			return err
		}

		if err := v.setupDirectoryStructure(); err != nil {
			return err
		}

		if err := v.saveConfig(); err != nil {
			return err
		}

	}
	return nil
}

func (v *Vault) generateKey() ([]byte, error) {
	bytes := make([]byte, 128)
	count, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	if count != 128 {
		return nil, errors.New("Not enough random bytes generated")
	}
	return bytes, nil
}

func (v *Vault) setupDirectoryStructure() error {
	os.MkdirAll(path.Join(v.Path, "config"), 0700)
	ioutil.WriteFile(path.Join(v.Path, "config", ".placeholder"), nil, 0600)
	os.MkdirAll(path.Join(v.Path, "keys"), 0700)
	ioutil.WriteFile(path.Join(v.Path, "keys", ".placeholder"), nil, 0600)
	return nil
}

func (v *Vault) Save(commitMsg string) error {
	return v.git.CommitAllChanges(commitMsg)
}

// seed the repository
func (v *Vault) Seed(creds *SshKeyRing) error {
	masterPassphrase, err := v.generateKey()
	if err != nil {
		return err
	}

	str, err := creds.EncryptAndBase64(masterPassphrase)
	if err != nil {
		log.Println("XXX ERROR", err)
		return err
	}

	return v.users.AddUser(v.email, creds.PublicKeyString(), str)
}
