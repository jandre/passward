package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/jandre/passward/passward"
	"github.com/segmentio/go-prompt"
)

func VaultSync(name string) {

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

	err = vault.Sync()

	if err != nil {
		fmt.Println("Unable to sync vault to remote store, did you call `passward vault set-remote`?")
		fmt.Println("Error is:", err)
		fmt.Println("If authentication fails, you may also need to add the following ssh public key to the remote git server:")
		fmt.Println("")
		fmt.Println(pw.Credentials.PublicKeyString())
		os.Exit(1)
	}

	fmt.Printf("Vault synced successfully: %s\n", vault.Name)
}
