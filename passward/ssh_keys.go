package passward

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/jandre/passward/util"
	"github.com/jandre/sshcrypt"
)

//
// Ssh keys management
//
type SshKeyRing struct {
	PublicKeyPath  string
	PrivateKeyPath string

	privateKey sshcrypt.PrivateKey
	publicKey  sshcrypt.PublicKey
}

func NewSshKeyRing(publicKeyPath string, privateKeyPath string, passphrase string) (*SshKeyRing, error) {

	ssh := SshKeyRing{PublicKeyPath: publicKeyPath, PrivateKeyPath: privateKeyPath}

	err := ssh.ParsePublicKey()
	if err != nil {
		return nil, err
	}

	if err = ssh.ParsePrivateKey(passphrase); err != nil {
		return nil, err
	}
	return &ssh, nil
}

//
// Get description string
//
func (s *SshKeyRing) GetDescription() string {
	return fmt.Sprintf("%s (Public), %s (Private)", s.PublicKeyPath, s.PrivateKeyPath)
}

// validates public key is ok and works with private key
func (s *SshKeyRing) ParsePublicKey() error {

	keyBytes, err := ioutil.ReadFile(s.PublicKeyPath)

	if err != nil {
		return err
	}

	ret, comment, opts, _, err := sshcrypt.ParseAuthorizedKey(keyBytes)

	if err != nil {
		return err
	}
	debug("read key with comment: %s, %s", comment, opts)
	s.publicKey = ret

	return nil
}

//
// ParsePrivateKey will parse a private key.
// Private keys in ssh are PEM-encoded blocks. Attempt to decode
//
// TODO: supports RSA only, add other support.
//
func (s *SshKeyRing) ParsePrivateKey(passphrase string) error {

	encryptedBytes, err := ioutil.ReadFile(s.PrivateKeyPath)

	if err != nil {
		return err
	}

	s.privateKey, err = sshcrypt.ParsePrivateKey(encryptedBytes, passphrase)

	if err != nil {
		return err

	}

	return nil
}

func GetSshKeyRingPath() string {
	home := os.Getenv("HOME")
	return path.Join(home, ".ssh")
}

//
// list all ssh keys in ~/.ssh
//
func DetectSshKeyRing(sshKeysPath string) []*SshKeyRing {
	keys := make([]*SshKeyRing, 0)
	if util.DirectoryExists(sshKeysPath) {

		files, err := filepath.Glob(path.Join(sshKeysPath, "*.pub"))
		if err != nil {
			// TODO: log this maybe?
			return keys
		}

		for _, pubKeyFile := range files {
			privateKeyFile := pubKeyFile[:len(pubKeyFile)-4]
			if util.FileExists(privateKeyFile) {
				key := &SshKeyRing{PublicKeyPath: pubKeyFile, PrivateKeyPath: privateKeyFile}
				keys = append(keys, key)
			}
		}
	}

	return keys
}
