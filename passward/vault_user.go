package passward

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"path"

	"github.com/jandre/passward/util"
	"github.com/jandre/sshcrypt"
)

//
// VaultUsers store the users in ~/.passward/vaults/<vault>/users/...
//
type VaultUsers struct {
	path  string // path to users/ directory for vault
	users map[string]*VaultUser
}

func NewVaultUsers(parentPath string) *VaultUsers {
	vaultUsersFolder := path.Join(parentPath, "users")
	result := VaultUsers{path: vaultUsersFolder, users: make(map[string]*VaultUser, 0)}
	return &result
}

// add a new vault user
func (vu *VaultUsers) AddUser(email string, publicKeyString string, masterPassphrase []byte) error {
	user, err := NewVaultUser(vu.path, email, publicKeyString)
	if err != nil {
		return err
	}

	if err := user.SetEncryptedMasterKey(masterPassphrase); err != nil {
		return err
	}

	if err := user.Save(); err != nil {
		return err
	}
	vu.users[email] = user
	return nil
}

func (vu *VaultUsers) Initialize() error {
	if !util.DirectoryExists(vu.path) {
		os.MkdirAll(vu.path, 0700)
		ioutil.WriteFile(path.Join(vu.path, ".placeholder"), nil, 0700)
	} else {

		files, _ := ioutil.ReadDir(vu.path)
		for _, name := range files {
			if name.Name() != ".placeholder" {
				user, err := ReadVaultUser(path.Join(vu.path, name.Name()))
				if err != nil {
					return err
				}
				vu.users[user.email] = user
			}

		}
	}
	return nil
}

func (vusers *VaultUsers) LookupByEmail(email string) *VaultUser {
	return vusers.users[email]
}

type VaultUser struct {
	path               string
	email              string
	publicKeyString    string
	encryptedMasterKey string
	publicKey          sshcrypt.PublicKey
}

func (vu *VaultUser) UnlockMasterKey(keyring *SshKeyRing) ([]byte, error) {
	// TODO: assert public keys match??
	return keyring.DecryptBase64(vu.encryptedMasterKey)
}

func (vu *VaultUser) Email() string {
	return vu.email
}
func (vu *VaultUser) PublicKey() string {
	return vu.publicKeyString
}

func (vu *VaultUser) Save() error {
	if !util.DirectoryExists(vu.path) {
		if err := os.MkdirAll(vu.path, 0700); err != nil {
			return err
		}
	}
	keyfile := vu.publicKeyFile()
	if err := ioutil.WriteFile(keyfile, []byte(vu.publicKeyString), 0600); err != nil {
		return err
	}

	encryptedMaster := vu.encryptedMasterFile()
	if err := ioutil.WriteFile(encryptedMaster, []byte(vu.encryptedMasterKey), 0600); err != nil {
		return err
	}
	return nil
}

func (vu *VaultUser) encryptedMasterFile() string {
	return path.Join(vu.path, "encrypted_master")
}

func (vu *VaultUser) publicKeyFile() string {
	return path.Join(vu.path, "key")
}

func (vu *VaultUser) GetEncryptedMasterKey() string {
	return vu.encryptedMasterKey
}

func (vu *VaultUser) SetEncryptedMasterKey(masterPassphrase []byte) error {
	cipherText, err := vu.publicKey.EncryptBytes(masterPassphrase)
	if err != nil {
		debug("failure to encrypt user master key: %s", err)
		return err
	}

	vu.encryptedMasterKey = base64.StdEncoding.EncodeToString(cipherText)
	return nil
}

func NewVaultUser(usersPath string, email string, publicKey string) (*VaultUser, error) {
	var user VaultUser
	var err error
	user.path = path.Join(usersPath, email)
	user.email = email
	user.publicKeyString = publicKey

	user.publicKey, _, _, _, err = sshcrypt.ParseAuthorizedKey([]byte(publicKey))

	if err != nil {
		debug("unable to parse public key:", publicKey, err)
		return nil, err
	}

	return &user, nil
}

func ReadVaultUser(pathToUser string) (*VaultUser, error) {
	var user VaultUser
	user.path = pathToUser
	user.email = path.Base(pathToUser)

	bytes, err := ioutil.ReadFile(user.publicKeyFile())

	if err != nil {
		debug("unable to parse public key %s", err)
		return nil, err
	}

	user.publicKeyString = string(bytes)
	user.publicKey, _, _, _, err = sshcrypt.ParseAuthorizedKey([]byte(user.publicKeyString))

	keyBytes, err := ioutil.ReadFile(user.encryptedMasterFile())

	if err != nil {
		debug("unable to parse master key %s", err)
		return nil, err
	}

	user.encryptedMasterKey = string(keyBytes)
	return &user, nil
}
