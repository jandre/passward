package passward

import (
	"os"
	"path"

	"github.com/jandre/passward/util"
	git2go "gopkg.in/libgit2/git2go.v22"
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
	Name     string
	Upstream string
	Path     string
	Entries  map[string]Entry
	repo     *git2go.Repository
	users    *VaultUsers
}

type VaultUsers struct {
	Path string
}

func NewVaultUsers(parentPath string) *VaultUsers {

	vaultUsersFolder = path.Join(parentPath, "users")

	result := VaultUsers{Path: vaultUsersFolder}

	return &result

}

type User struct {
	Name  string
	Keys  []string
	Email string
}

//
// Create a new vault.
//
func NewVault(config *Passward, name string) (*Vault, error) {
	dst := path.Join(config.Path, name)

	result := Vault{Name: name, Path: dst, users: NewVaultUsers(dst)}

	if err := result.Initialize(); err != nil {
		return nil, err
	}
	return &result, nil
}

//
// Initialize the new vault by performing a git init at path, etc.
//
func (v *Vault) Initialize() error {

	// it's already setup
	if util.DirectoryExists(v.Path) {
		return nil
	}

	if repo, err := git2go.InitRepository(v.Path, false); err != nil {
		return err
	} else {
		v.repo = repo
	}

	return v.setupDirectoryStructure()
}

func (v *Vault) setupDirectoryStructure() error {
	os.MkdirAll(path.Join(v.Path, "users"), 0700)
	os.MkdirAll(path.Join(v.Path, "config"), 0700)
	os.MkdirAll(path.Join(v.Path, "keys"), 0700)
	return nil
}

func (v *Vault) LoadVault()    {}
func (v *Vault) AddEntry()     {}
func (v *Vault) DeleteEntry()  {}
func (v *Vault) AddUser()      {}
func (v *Vault) RemoveUser()   {}
func (v *Vault) SetMasterKey() {}
