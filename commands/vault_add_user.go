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

	var vault *passward.Vault

	if name != "" {
		vault = pw.GetVault(name)

		if vault == nil {
			log.Fatal("Vault not found: " + name)
		}
	} else {
		vault = pw.GetSelectedVault()
		if vault == nil {
			log.Fatal("No vault found; you need to run `passward vault use name` to select a vault or `passward vault new <name>` to create one.")
		}
	}

	passphrase := prompt.PasswordMasked("Enter your passphrase to unlock your keys (empty for none)")
	if err := pw.Unlock(passphrase); err != nil {
		log.Fatal("Invalid passphrase.", err)
	}

	log.Println("Please enter the public key (e.g. the contents of ~/.ssh/id_rsa.pub).")
	publicKey := prompt.StringRequired("Enter key")

	_, err = vault.AddUser(email, publicKey)
	if err != nil {
		log.Fatal("Unable to set add user: ", err)
	}

	fmt.Printf("User for `%s` successfully added: %s.\n", vault.Name, email)
	fmt.Println("You will want to ensure that the ssh key has permission to the remote repository.")
}
