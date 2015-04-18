package passward

import (
	"errors"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/jandre/passward/util"
)

type Passward struct {
	Path          string `toml:"-"`
	Credentials   *Credentials
	SelectedVault string
	vaults        map[string]*Vault
	selectedVault *Vault
}

func (pw *Passward) GetSelectedVault() *Vault {
	return pw.selectedVault
}

func (pw *Passward) UseVault(name string) error {
	vault := pw.vaults[name]
	if vault == nil {
		return errors.New("vault not found:" + name)
	}
	pw.SelectedVault = name
	pw.selectedVault = vault
	return nil
}

//
// GetVaults vaults in ~/.passward/vaults
//
func (pw *Passward) GetVaults() map[string]*Vault {
	return pw.vaults
}

//
// Get a vault with name = `name`
//
func (pw *Passward) GetVault(name string) *Vault {
	return pw.vaults[name]
}

func (pw *Passward) SetCredentials(creds *Credentials) {
	pw.Credentials = creds
}

func (pw *Passward) GetCredentials() *Credentials {
	if pw.Credentials == nil {
		panic("credentials not set")
	}
	return pw.Credentials
}

//
// Unlock a passward by supplying credentials
//
func (pw *Passward) Unlock(passphrase string) error {
	if pw.Credentials == nil {
		panic("no credentials set!")
	}
	return pw.Credentials.Unlock(passphrase)
}

func (pw *Passward) FetchVault(url string, name string) (*Vault, error) {
	if name == "" {
		name = detectGitName(url)
	}

	if pw.vaults[name] != nil {
		return nil, errors.New("Vault " + name + " already exists!")
	}

	tmpDir := path.Join(pw.Path, "vaults", name)

	creds := pw.GetCredentials()
	// make a tmpdir
	git := NewGit(tmpDir, creds)

	debug("cloning to ", tmpDir)

	err := git.Clone(url)

	if err != nil {
		return nil, err
	}

	vault, err := ReadVault(pw.vaultPath(), name, pw.GetCredentials())

	if err != nil {
		return nil, err
	}

	if vault != nil {
		pw.vaults[name] = vault
	} else {
		return nil, errors.New("Unable to fetch, no vault found:" + name)
	}

	return vault, nil
}

//
// Add a vault to the ~/.passward/vaults
//
func (pw *Passward) AddVault(name string) error {
	// check to see if there is a vault with the name already
	if pw.vaults[name] != nil {
		return errors.New("Vault " + name + " already exists!")
	}

	creds := pw.GetCredentials()

	if vault, err := NewVault(pw.vaultPath(), name, creds); err != nil {
		return err
	} else {

		if !creds.IsUnlocked() {
			return errors.New("Credentials must be unlocked.")
		}

		if err = vault.Initialize(); err != nil {
			return err
		}

		if err = vault.Seed(); err != nil {
			return err
		}

		if err = vault.Save("New vault created."); err != nil {
			return err
		}
		pw.vaults[name] = vault
	}

	return nil
}

func (c *Passward) vaultPath() string {
	return path.Join(c.Path, "vaults")
}

func (c *Passward) configPath() string {
	return path.Join(c.Path, "config.toml")
}

func DetectPasswardPath() string {
	// first looking at PASSWARD_HOME, then ~/.passward, then /opt/passward
	if p := os.Getenv("PASSWARD_HOME"); p != "" {
		return p
	}

	home := os.Getenv("HOME")
	if home == "" {
		log.Println("no home detected! using global path /opt/passward")
		return "/opt/passward"
	}
	return filepath.Join(home, ".passward")
}

func NewPassward(directory string) (*Passward, error) {
	var conf Passward
	if directory == "" {
		directory = DetectPasswardPath()
	}

	directory, err := filepath.Abs(directory)

	if err != nil {
		return nil, err
	}

	if util.DirectoryExists(directory) {
		return nil, errors.New("passward home already exists:" + directory)
	}

	conf.vaults = make(map[string]*Vault, 0)
	conf.Path = directory

	os.MkdirAll(conf.vaultPath(), 0700)

	return &conf, nil
}

//
// ReadPassward() will parse a passward config at `configPath`
//
func ReadPassward(directory string) (*Passward, error) {
	var pw Passward
	directory, err := filepath.Abs(directory)

	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(directory, "config.toml")
	_, err = toml.DecodeFile(configPath, &pw)

	if err != nil {
		return nil, err
	}
	pw.Path = directory // in case it was moved
	pw.vaults, err = ReadAllVaults(pw.vaultPath(), pw.GetCredentials())
	if err != nil {
		return nil, err
	}

	if pw.SelectedVault != "" {
		pw.UseVault(pw.SelectedVault)
	}

	return &pw, nil
}

//
// Save saves the config file to ~/passward/config.toml
//
func (c *Passward) Save() error {

	if !util.DirectoryExists(c.Path) {
		if err := os.MkdirAll(c.Path, 0700); err != nil {
			return err
		}
	}

	file, err := os.Create(c.configPath())
	if err != nil {
		return err
	}
	defer file.Close()
	if err := toml.NewEncoder(file).Encode(c); err != nil {
		return err
	}

	return nil
}
