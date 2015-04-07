package passward

import (
	"crypto/dsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/jandre/passward/util"
	"golang.org/x/crypto/ssh"
)

type rsaPrivateKey struct {
	*rsa.PrivateKey
}

type dsaPrivateKey struct {
	*dsa.PrivateKey
}

type rsaPublicKey struct {
	*rsa.PublicKey
}

func (pk *rsaPublicKey) Type() string {
	return "ssh-rsa"
}

func (pk *rsaPublicKey) Marshal() []byte {
	return nil
}

func (pk *rsaPublicKey) Verify(data []byte, sig *ssh.Signature) error {
	return nil
}

func (pk *rsaPublicKey) toPEM() (string, error) {
	key, err := x509.MarshalPKIXPublicKey(pk)
	if err != nil {
		return "", err
	}
	block := pem.Block{Type: "BEGIN PUBLIC KEY", Bytes: key}
	return string(pem.EncodeToMemory(&block)), nil
}

type dsaPublicKey struct {
	*dsa.PublicKey
}

//
// Ssh keys management
//
type SshKeys struct {
	PublicKeyPath  string
	PrivateKeyPath string

	privateKey interface{}
	publicKey  ssh.PublicKey

	privateKeyPem string
	publicKeyPem  string
}

func (s *SshKeys) GetDescription() string {
	return fmt.Sprintf("%s (Public), %s (Private)", s.PublicKeyPath, s.PrivateKeyPath)
}

func (s *SshKeys) ParseKeys(passphrase string) error {

	return nil
}

// validates public key is ok and works with private key
func (s *SshKeys) ParsePublicKey() error {

	keyBytes, err := ioutil.ReadFile(s.PublicKeyPath)

	if err != nil {
		return err
	}

	k, comment, opts, _, err := ssh.ParseAuthorizedKey(keyBytes)

	s.publicKey = k

	debug("read key with comment: %s, %s", comment, opts)

	return err
}

//
// ParsePrivateKey will parse a private key.
// Private keys in ssh are PEM-encoded blocks. Attempt to decode
//
// TODO: supports RSA only, add other support.
//
func (s *SshKeys) ParsePrivateKey(passphrase string) error {

	encryptedBytes, err := ioutil.ReadFile(s.PrivateKeyPath)

	s.privateKeyPem = string(encryptedBytes)

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
