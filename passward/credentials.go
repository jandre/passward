package passward

import "io/ioutil"

type Credentials struct {
	keyring        *SshKeyRing `toml:"-"`
	keyPassphrase  string      `toml:"-"`
	Name           string
	Email          string
	PublicKeyPath  string
	PrivateKeyPath string
}

func (creds *Credentials) PublicKeyString() string {

	if creds.keyring != nil {
		return creds.keyring.PublicKeyString()
	} else {
		bytes, err := ioutil.ReadFile(creds.PublicKeyPath)
		if err != nil {
			return ""
		}
		return string(bytes)
	}

}

func (creds *Credentials) Passphrase() string {
	return creds.keyPassphrase
}

func (creds *Credentials) Lock() {
	creds.keyring = nil
	creds.keyPassphrase = ""
}

func (creds *Credentials) GetKeys() *SshKeyRing {
	if creds.keyring == nil {
		panic("Credentials have not been unlocked")
	}
	return creds.keyring
}

func (creds *Credentials) Unlock(passphrase string) error {
	var err error

	creds.keyPassphrase = passphrase
	if creds.keyring != nil {
		debug("already unlocked")
		return nil
	}
	creds.keyring, err = NewSshKeyRing(creds.PublicKeyPath, creds.PrivateKeyPath, passphrase)
	return err
}

func (creds *Credentials) IsUnlocked() bool {
	return creds.keyring != nil
}
