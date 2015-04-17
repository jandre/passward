package passward

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/jandre/passward/util"
)

type Entry struct {
	name            string
	path            string
	encryptedValues map[string]string
}

func NewEntry(parentDir, name string) *Entry {
	entry := Entry{
		name:            name,
		path:            path.Join(parentDir, name),
		encryptedValues: make(map[string]string),
	}

	return &entry
}

func ReadEntry(parentDir, name string) (*Entry, error) {
	entry := NewEntry(parentDir, name)
	files, err := ioutil.ReadDir(entry.path)

	if err != nil {
		return nil, err
	}

	for _, file := range files {
		filename := file.Name()
		bytes, err := ioutil.ReadFile(path.Join(entry.path, filename))
		if err != nil {
			return nil, err
		}

		entry.encryptedValues[filename] = string(bytes)
	}

	return entry, nil
}

func (e *Entry) Name() string {
	return e.name
}

func (e *Entry) Set(key string, val string, encryptionKey []byte) error {
	cryptKey := string(encryptionKey)
	encryptedVal, err := EncryptAndBase64String(cryptKey, val)
	if err != nil {
		return err
	}

	e.encryptedValues[key] = encryptedVal
	return nil
}

func (e *Entry) Save() error {
	if !util.DirectoryExists(e.path) {
		os.MkdirAll(e.path, 0700)
	}
	for key, val := range e.encryptedValues {
		file := path.Join(e.path, key)
		err := ioutil.WriteFile(file, []byte(val), 0600)
		if err != nil {
			return err
		}
	}
	return nil
}

type VaultEntries struct {
	entries map[string]*Entry
	path    string
}

func NewVaultEntries(parentDir string) *VaultEntries {
	ve := VaultEntries{
		entries: make(map[string]*Entry, 0),
		path:    path.Join(parentDir, "keys"),
	}
	return &ve
}

func (ve *VaultEntries) Path() string {
	return ve.path
}

func (ve *VaultEntries) Initialize() error {
	if !util.DirectoryExists(ve.path) {
		if err := os.MkdirAll(ve.path, 0700); err != nil {
			return err
		}
		if err := ioutil.WriteFile(path.Join(ve.path, ".placeholder"), nil, 0600); err != nil {
			return err
		}
	}

	files, err := ioutil.ReadDir(ve.path)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.Name() != ".placeholder" {
			entry, err := ReadEntry(ve.Path(), file.Name())

			if err != nil {
				debug("unable to load entry", err)
				return err
			}

			ve.entries[entry.Name()] = entry
		}
	}

	return nil
}

func (ve *VaultEntries) Add(name string, key string, val string, encryptionKey []byte) error {
	if ve.entries[name] == nil {
		ve.entries[name] = NewEntry(ve.Path(), name)
	}

	return ve.entries[name].Set(key, val, encryptionKey)
}

func (ve *VaultEntries) Save() error {
	for _, entry := range ve.entries {
		err := entry.Save()
		if err != nil {
			return err
		}
	}
	return nil
}
