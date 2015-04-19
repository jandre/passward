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

const KEYSIZE = 128

type Vault struct {
	Name        string
	Description string
	Path        string `toml:"-"`

	// secret entries
	entries *VaultEntries `toml:"-"`

	// users in vault
	users       *VaultUsers  `toml:"-"`
	credentials *Credentials `toml:"-"`
	git         *Git         `toml:"-"`
}

func (v *Vault) RemoveUser(email string) error {
	panic("not implemented")
	return nil
}

func (v *Vault) GetUserByEmail(email string) *VaultUser {
	return v.users.LookupByEmail(email)
}

//
// HasRemote returns true if the vault has a remote set.
//
func (v *Vault) HasRemote() bool {
	return v.git.HasRemote()
}

func (v *Vault) AddUser(email string, publicKey string) (*VaultUser, error) {
	masterKey, err := v.unlockMasterKey()
	if err != nil {
		debug("could not add user - vault is not unlocked")
		return nil, err
	}

	if err := v.users.AddUser(email, publicKey, masterKey); err != nil {
		return nil, err
	}

	user := v.users.LookupByEmail(email)

	v.Save("Added user: " + email)

	return user, nil
}

func (v *Vault) Users() map[string]*VaultUser {
	return v.users.users
}

func (v *Vault) Entries() map[string]*Entry {
	return v.entries.entries
}

func (v *Vault) unlockMasterKey() ([]byte, error) {

	keys := v.credentials.GetKeys()

	if keys == nil {
		return nil, errors.New("No keys found - did you call passward.Unlock()?")
	}

	user := v.users.LookupByEmail(v.credentials.Email)

	if user == nil {
		return nil, errors.New("No vault user found to unlock passphrase")
	}

	return user.UnlockMasterKey(keys)
}

func (v *Vault) RevealEntry(name string) (secrets map[string]string, err error) {
	key, err := v.unlockMasterKey()
	if err != nil {
		return nil, err
	}

	entry := v.entries.Get(name)
	if entry == nil {
		return nil, errors.New("No entry found:" + name)
	}

	return entry.RevealAll(key)
}

func (v *Vault) AddEntry(name string, user string, passphrase string, desc string) error {
	key, err := v.unlockMasterKey()
	if err != nil {
		return err
	}
	v.entries.Add(name, "username", user, key)
	v.entries.Add(name, "passphrase", passphrase, key)
	v.entries.Add(name, "description", desc, key)
	if err := v.entries.Save(); err != nil {
		return err
	}
	return v.Save("New entry: " + name)
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
	vault.git = NewGit(dst, creds)
	vault.entries = NewVaultEntries(dst)
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
		entries:     NewVaultEntries(dst),
		credentials: creds,
		git:         NewGit(dst, creds),
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
		if err = v.entries.Initialize(); err != nil {
			return err
		}
	} else {

		debug("initializing vault: %s", v.Path)

		if err = v.git.Initialize(); err != nil {
			debug("error initializing git: %s", err)
			return err
		}

		if err = v.entries.Initialize(); err != nil {
			debug("error initializing entries: %s", err)
			return err
		}

		if err = v.users.Initialize(); err != nil {
			debug("error initializing users: %s", err)
			return err
		}
		if err := v.saveConfig(); err != nil {
			return err
		}
	}
	return nil
}

func (v *Vault) generateKey() ([]byte, error) {

	bytes := make([]byte, KEYSIZE)
	count, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	if count != KEYSIZE {
		return nil, errors.New("Not enough random bytes generated")
	}
	return bytes, nil
}

func (v *Vault) Save(commitMsg string) error {
	return v.git.CommitAllChanges(commitMsg)
}

func (v *Vault) SetRemote(remote string) error {
	return v.git.SetRemote(remote)
}

func (v *Vault) Sync() error {
	return v.git.Push()

	// TODO: also pull
}

// seed the repository
func (v *Vault) Seed() error {
	masterPassphrase, err := v.generateKey()
	if err != nil {
		return err
	}

	return v.users.AddUser(v.credentials.Email, v.credentials.GetKeys().PublicKeyString(), masterPassphrase)
}
