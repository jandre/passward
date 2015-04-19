package commands

import (
	"fmt"
	"log"

	"github.com/jandre/passward/passward"
	prompt "github.com/segmentio/go-prompt"
)

func VaultAddUser(name string, email string) {

	passwardPath := passward.DetectPasswardPath()

	pw, err := passward.ReadPassward(passwardPath)

	if err != nil {
		log.Fatal("There was a problem loading the configuration. Did you run `passward setup?`", err)
	}

	vault := chooseVault(pw, name)

	passphrase := prompt.PasswordMasked("Enter your passphrase to unlock your keys (empty for none)")
	if err := pw.Unlock(passphrase); err != nil {
		log.Fatal("Invalid passphrase.", err)
	}

	log.Println("Please enter the public key (e.g. the contents of ~/.ssh/id_rsa.pub).")
	publicKey := prompt.StringRequired("Enter key")

	_, err = vault.AddUser(email, publicKey)
	if err != nil {
		log.Fatal("Unable to add user: ", err)
	}

	fmt.Printf("User `%s` successfully saved to vault: %s.\n", email, vault.Name)
	if vault.HasRemote() {
		fmt.Println("1. Sync your changes by running `passward vault sync`.")
		fmt.Println("2. You will want to ensure that the ssh key has permission to the remote repository.")
	}
}
