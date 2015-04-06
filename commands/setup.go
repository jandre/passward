package commands

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/jandre/passward/passward"
	"github.com/jandre/passward/util"
	prompt "github.com/segmentio/go-prompt"
)

type SshKeys struct {
	PublicKeyPath  string
	PrivateKeyPath string
}

func (s *SshKeys) GetDescription() string {
	return fmt.Sprintf("%s (Public), %s (Private)", s.PublicKeyPath, s.PrivateKeyPath)
}

func makeSshKeyDescriptions(keys []*SshKeys) []string {
	result := make([]string, len(keys))
	for i, k := range keys {

		result[i] = k.GetDescription()
	}
	return result

}

func getSshKeysPath() string {
	home := os.Getenv("HOME")
	return path.Join(home, ".ssh")
}

//
// list all ssh keys in ~/.ssh
//
func detectSshKeys(sshKeysPath string) []*SshKeys {
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

func ChooseSshKeys() *SshKeys {
	sshKeysPath := getSshKeysPath()
	sshKeys := detectSshKeys(sshKeysPath)

	if sshKeys != nil && len(sshKeys) > 0 {
		fmt.Printf("passward uses your SSH keys for encryption. We've detected the following keys, choose the ones you want to use: \n")
		sshKeyDescriptions := makeSshKeyDescriptions(sshKeys)
		sshKeyDescriptions = append(sshKeyDescriptions, "None of these, generate new keys for me.")
		id := prompt.Choose("Select keys to use", sshKeyDescriptions)
		if id != len(sshKeys) {
			return sshKeys[id]
		} else {
			// TODO: have the app run ssh-keygen all on its lonesome.
			fmt.Println("Please run ssh-keygen to generate the keys.")
			os.Exit(1)
		}
	} else {
		fmt.Printf("No ssh keys detected in %s!\n", sshKeysPath)
		fmt.Println("Please run ssh-keygen to generate the keys.")
		os.Exit(1)
	}
	return nil
}

//
// Setup a new passward installation
//
func Setup() {
	passwardPath := passward.DetectPasswardPath()

	if util.DirectoryExists(passwardPath) {
		fmt.Println("Oh no! We already detected a passward installation at: ", passwardPath, ".")
		fmt.Println("Please remove this directory, or set environment variable PASSWARD_HOME=<path> to use a different path.")
	}

	fmt.Println("Hello! We'll be installing passward here: ", passwardPath)
	fmt.Println("(If you don't want it here, please export PASSWARD_HOME=<blah> and re-run `passward setup`.")

	fmt.Println("--")

	sshKeys := ChooseSshKeys()

	fmt.Println("Great! We'll be using the keys at: ", sshKeys.GetDescription())
	prompt.PasswordMasked("Please enter the passphrase for these keys")

	//prompt.Confirm("We'll be installing .passward at %s, ? ", passwardPath)

	//	installation := passward.NewPassward()

}
