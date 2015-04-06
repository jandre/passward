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
		fmt.Printf("Wonderful. We've detected the following keypairs, choose the ones you want to use: \n")
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

func chooseAuthMethod(cfg *passward.Passward) {

	authMethods := []string{
		"Import SSH keys - Use your ssh keys to encrypt your password vaults.",
		"Generate custom keys - Generate custom public keypair.",
	}

	fmt.Println("Passward uses public key encryption to store secrets. You can use existing keys from SSH, or generate new ones.")
	chosen := prompt.Choose("Select your authentication method", authMethods)

	if chosen == 0 {
		setupSshAuth(cfg)
	} else {
		log.Fatal("Sorry, this is not yet supported!")
	}
}

func setupSshAuth(cfg *passward.Passward) {
	sshKeys := ChooseSshKeys()

	fmt.Println()
	fmt.Println("Great! We'll be using the keypair: ", sshKeys.GetDescription())

	found := false
	attempts := 1

	for !found {
		var passphrase string
		if attempts > 1 {
			passphrase = prompt.PasswordMasked(
				fmt.Sprintf("(attempt %d/3) Please enter the passphrase for the private key", attempts))
		} else {
			passphrase = prompt.PasswordMasked("Please enter the passphrase for the private key, or just hit enter if there is no password")
		}
		err := sshKeys.ParsePrivateKey(passphrase)
		if err == nil {
			break
		}

		if attempts >= 3 {
			log.Println("Unable to decrypt private key due to:", err)
			log.Fatal("Looks like there was a major problem decrypting or using your private key.")

		}
		attempts++
	}

	err := sshKeys.ParsePublicKey()

	if err != nil {
		log.Fatal("Unable to parse public key:", err)
	}
	fmt.Println("Great! We've loaded the keys and verified everything is great.")

	// convert public key to PEM forma

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

	fmt.Println("Hello! Welcome to passward setup. We'll be installing passward here: ", passwardPath)
	fmt.Println("(If you don't want the data files installed here, please export PASSWARD_HOME=<blah> and re-run `passward setup`.)")

	fmt.Println("")

	email := prompt.StringRequired("Please enter your email. Your email address is also your passward username")

	cfg, err := passward.NewPassward(email, passwardPath)

	if err != nil {
		log.Fatal("Unable to save passward config. You can try to re-run `passward setup`.", err)
	}

	chooseAuthMethod(cfg)

	err = cfg.Save()

	if err != nil {
		log.Fatal("Unable to save passward config. You can try to re-run `passward setup`.", err)
	}

	//prompt.Confirm("We'll be installing .passward at %s, ? ", passwardPath)
	//	installation := passward.NewPassward()

}
