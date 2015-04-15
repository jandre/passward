package passward

import (
	"crypto/rand"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/jandre/passward/util"
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
	credentials    *Credentials     `toml:"-"`
	git            *Git             `toml:"-"`
}

func (v *Vault) Users() map[string]*VaultUser {
	return v.users.users
}

func (v *Vault) unlockMasterPassphrase() ([]byte, error) {

	// keys := v.credentials.GetKeys()
	// user := v.users.LookupByEmail(v.credentials.Email)

	// if user == nil {
	// return nil, errors.New("No vault user found to unlock passphrase")
	// }

	// XX: TODO
	return nil, nil
}

func (v *Vault) AddEntry(name string, user string, passphrase string, desc string) error {
	// XXX: todo
	return nil
}

func ReadAllVaults(vaultPath string, creds *Credentials) (map[string]*Vault, error) {
	vaults := make(map[string]*Vault, 0)

	files, err := ioutil.ReadDir(vaultPath)

	if err != nil {
		return nil, err
	}

	for _, name := range files {
		vault, err := ReadVault(vaultPath, name.Name(), creds)
		if err != nil {
			return nil, err
		}

		if vault != nil {
			vaults[name.Name()] = vault
		}
	}

	return vaults, nil
}

func ReadVault(vaultPath string, name string, creds *Credentials) (*Vault, error) {
	dst := path.Join(vaultPath, name)
	configPath := filepath.Join(dst, "config.toml")

	var vault Vault
	_, err := toml.DecodeFile(configPath, &vault)

	if err != nil {
		return nil, err
	}

	vault.Path = dst // in case it was moved
	vault.users = NewVaultUsers(dst)
	vault.credentials = creds
	vault.git = NewGit(dst, creds.Name, creds.Email)
	vault.Initialize()
	return &vault, nil
}

//
// Create a new vault.
//
func NewVault(vaultPath string, name string, creds *Credentials) (*Vault, error) {
	dst := path.Join(vaultPath, name)

	result := Vault{Name: name,
		Path:        dst,
		users:       NewVaultUsers(dst),
		credentials: creds,
		git:         NewGit(dst, creds.Name, creds.Email),
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
func (v *Vault) Seed() error {
	masterPassphrase, err := v.generateKey()
	if err != nil {
		return err
	}

	return v.users.AddUser(v.credentials.Email, v.credentials.GetKeys().PublicKeyString(), masterPassphrase)
}
