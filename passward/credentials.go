package passward

type Credentials struct {
	keyring        *SshKeyRing `toml:"-"`
	keyPassphrase  string      `toml:"-"`
	Name           string
	Email          string
	PublicKeyPath  string
	PrivateKeyPath string
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
