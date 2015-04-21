package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/jandre/passward/passward"
	prompt "github.com/segmentio/go-prompt"
)

func VaultRemoveUser(name string, email string) {

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

	user := vault.GetUserByEmail(email)

	if user == nil {
		fmt.Printf("User %s could not be found in vault %s.\n", email, vault.Name)
		os.Exit(1)
	}

	prompt.Confirm(fmt.Sprintf("Are you sure you want to remove the user %s?", email))

	err = vault.RemoveUser(email)
	if err != nil {
		log.Fatal("Unable to remove user: ", err)
	}

}
