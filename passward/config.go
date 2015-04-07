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
	Vaults     []Vault
	Email      string
	Path       string `toml:"-"`
	PrivateKey string
	PublicKey  string
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

	conf.Path = directory
	conf.Email = email

	return &conf, nil
}

//
// ReadPassward() will parse a passward config at `configPath`
//
func ReadPassward(directory string) (*Passward, error) {
	var conf Passward
	directory, err := filepath.Abs(directory)

	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(directory, "config.toml")
	_, err = toml.DecodeFile(configPath, &conf)

	conf.Path = directory // in case it was moved
	if err != nil {
		return nil, err
	}

	return &conf, nil
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
	return toml.NewEncoder(file).Encode(c)
}
