package passward

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/jandre/passward/util"
)

//
// Ssh keys management
//
type SshKeys struct {
	PublicKeyPath  string
	PrivateKeyPath string

	privateKey interface{}
}

func (s *SshKeys) GetDescription() string {
	return fmt.Sprintf("%s (Public), %s (Private)", s.PublicKeyPath, s.PrivateKeyPath)
}

func (s *SshKeys) ParseKeys(passphrase string) error {

	return nil
}

// validates public key is ok and works with private key
func (s *SshKeys) ParsePublicKey() error {

	return nil
}

//
// ParsePrivateKey will parse a private key.
// Private keys in ssh are PEM-encoded blocks. Attempt to decode
//
// TODO: supports RSA only, add other support.
//
func (s *SshKeys) ParsePrivateKey(passphrase string) error {

	encryptedBytes, err := ioutil.ReadFile(s.PrivateKeyPath)

	if err != nil {
		return err
	}

	decodedPEM, _ := pem.Decode(encryptedBytes)

	var privateBytes []byte
	if passphrase != "" {

		privateBytes, err = x509.DecryptPEMBlock(decodedPEM, []byte(passphrase))

		if err != nil {
			return err
		}
	} else {
		privateBytes = decodedPEM.Bytes
	}

	private, err := x509.ParsePKCS1PrivateKey(privateBytes)
	///	private, err := ssh.ParsePrivateKey(privateBytes)

	if err != nil {
		return err

	}

	s.privateKey = private

	return nil
}

func GetSshKeysPath() string {
	home := os.Getenv("HOME")
	return path.Join(home, ".ssh")
}

//
// list all ssh keys in ~/.ssh
//
func DetectSshKeys(sshKeysPath string) []*SshKeys {
	keys := make([]*SshKeys, 0)
	if util.DirectoryExists(sshKeysPath) {

		files, err := filepath.Glob(path.Join(sshKeysPath, "*.pub"))
		if err != nil {
			// TODO: log this maybe?
			return keys
		}

		for _, pubKeyFile := range files {
			privateKeyFile := pubKeyFile[:len(pubKeyFile)-4]
			if util.FileExists(privateKeyFile) {
				key := &SshKeys{PublicKeyPath: pubKeyFile, PrivateKeyPath: privateKeyFile}
				keys = append(keys, key)
			}
		}
	}

	return keys
}
