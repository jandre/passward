package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/jandre/passward/passward"
	"github.com/jandre/passward/util"
	prompt "github.com/segmentio/go-prompt"
)

func makeSshKeyDescriptions(keys []*passward.SshKeys) []string {
	result := make([]string, len(keys))
	for i, k := range keys {

		result[i] = k.GetDescription()
	}
	return result
}

func ChooseSshKeys() *passward.SshKeys {
	sshKeysPath := passward.GetSshKeysPath()
	sshKeys := passward.DetectSshKeys(sshKeysPath)

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

	fmt.Println("")

	sshKeys := ChooseSshKeys()

	fmt.Println()
	fmt.Println("Great! We'll be using the keys at: ", sshKeys.GetDescription())

	found := false
	attempts := 1

	for !found {
		var passphrase string
		if attempts > 1 {
			passphrase = prompt.PasswordMasked(
				fmt.Sprintf("(attempt %d/3) Please enter the passphrase for the private key", attempts))
		} else {
			passphrase = prompt.PasswordMasked("Please enter the passphrase for the private key")
		}
		err := sshKeys.ParsePrivateKey(passphrase)
		if err == nil {
			break
		}

		if attempts >= 3 {
			log.Fatal("Unable to decrypt private key due to:", err)

		}
		attempts++
	}

	fmt.Println("Great! We've loaded the key.")

	//prompt.Confirm("We'll be installing .passward at %s, ? ", passwardPath)
	//	installation := passward.NewPassward()

}
