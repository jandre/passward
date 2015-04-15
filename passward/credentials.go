package passward

type Credentials struct {
	keyring        *SshKeyRing `toml:"-"`
	Name           string
	Email          string
	PublicKeyPath  string
	PrivateKeyPath string
}

func (creds *Credentials) GetKeys() *SshKeyRing {
	if creds.keyring == nil {
		panic("Credentials have not been unlocked")
	}
	return creds.keyring
}

func (creds *Credentials) Unlock(passphrase string) error {
	var err error
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
