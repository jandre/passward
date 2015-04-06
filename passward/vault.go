package passward

type EncryptedValue struct {
	Path string
}

type Entry struct {
	Name   string
	Path   string
	Values map[string]EncryptedValue
}

type Vault struct {
	Name     string
	Upstream string
	Path     string
	Entries  map[string]Entry
}

type User struct {
	Name  string
	Keys  []string
	Email string
}

func (v *Vault) NewVault()     {}
func (v *Vault) LoadVault()    {}
func (v *Vault) AddEntry()     {}
func (v *Vault) DeleteEntry()  {}
func (v *Vault) AddUser()      {}
func (v *Vault) RemoveUser()   {}
func (v *Vault) SetMasterKey() {}
