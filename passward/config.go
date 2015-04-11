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
	Vaults     map[string]*Vault
	Email      string
	Path       string `toml:"-"`
	PrivateKey string
	PublicKey  string
}

func (pw *Passward) AddVault(name string) error {
	// check to see if there is a vault with the name already
	if pw.Vaults[name] != nil {
		return errors.New("Vault " + name + " already exists!")
	}

	if vault, err := NewVault(pw.vaultPath(), name, pw.Email, pw.Email); err != nil {
		return err
	} else {
		if err = vault.Initialize(); err != nil {
			return err
		}
		vault.Save("New vault created.")
		pw.Vaults[name] = vault
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

//
//
//
func NewPassward(email string, directory string) (*Passward, error) {
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

	conf.Vaults = make(map[string]*Vault, 0)
	conf.Path = directory
	conf.Email = email

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
	pw.Vaults, err = ReadAllVaults(pw.vaultPath())
	if err != nil {
		return nil, err
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
