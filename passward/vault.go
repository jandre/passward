package passward

import (
	"io/ioutil"
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
	Name           string
	RemoteUpstream string
	Path           string             `toml:"-"`
	Entries        map[string]Entry   `toml:"-"`
	repo           *git2go.Repository `toml:"-"`
	users          *VaultUsers        `toml:"-"`
}

type VaultUsers struct {
	Path string
}

func NewVaultUsers(parentPath string) *VaultUsers {
	vaultUsersFolder := path.Join(parentPath, "users")
	result := VaultUsers{Path: vaultUsersFolder}
	return &result
}

func (vu *VaultUsers) Initialize() {
	if !util.DirectoryExists(vu.Path) {
		os.MkdirAll(vu.Path, 0700)
	}

	// TODO: read users
}

type User struct {
	Name  string
	Keys  []string
	Email string
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

	result := Vault{Name: name, Path: dst, users: NewVaultUsers(dst)}
	result.users.Initialize()
	return &result, nil
}

//
// Create a new vault.
//
func NewVault(vaultPath string, name string) (*Vault, error) {
	dst := path.Join(vaultPath, name)

	result := Vault{Name: name, Path: dst, users: NewVaultUsers(dst)}

	result.users.Initialize()

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

	os.MkdirAll(path.Join(v.Path, "config"), 0700)
	os.MkdirAll(path.Join(v.Path, "keys"), 0700)
	return nil
}

func (v *Vault) Save(pw *Passward, commitMsg string) error {
	// sig := &git2go.Signature{
	// Name:  pw.Email,
	// Email: pw.Email,
	// }

	// idx, err := v.repo.Index()

	// if err != nil {
	// return err
	// }

	return nil
	// err = idx.AddByPath("README")
	// treeId, err := idx.WriteTree()
	// tree, err := repo.LookupTree(treeId)
	// commitId, err := repo.CreateCommit("HEAD", sig, sig, commitMsg, tree)

	// return commitId, treeId
}

func (v *Vault) LoadVault()    {}
func (v *Vault) AddEntry()     {}
func (v *Vault) DeleteEntry()  {}
func (v *Vault) AddUser()      {}
func (v *Vault) RemoveUser()   {}
func (v *Vault) SetMasterKey() {}
