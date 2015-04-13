package commands

import (
	"fmt"
	"log"

	"github.com/jandre/passward/passward"
	"github.com/segmentio/go-prompt"
)

func VaultNew(name string) {

	passwardPath := passward.DetectPasswardPath()

	pw, err := passward.ReadPassward(passwardPath)

	if err != nil {
		log.Fatal("There was a problem loading the configuration. Did you run `passward setup?`", err)
	}

	passphrase := prompt.PasswordMasked("Enter your passphrase to unlock your keys (empty for none)")

	if err := pw.Unlock(passphrase); err != nil {
		log.Fatal("Invalid passphrase.", err)
	}

	if err = pw.AddVault(name); err != nil {
		log.Fatal("Error creating vault: ", err)
	}

	fmt.Println("Successfully created new vault: ", name)

}
