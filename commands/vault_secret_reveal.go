package commands

import (
	"fmt"
	"log"

	"github.com/jandre/passward/passward"
	"github.com/segmentio/go-prompt"
)

func VaultSecretReveal(name string, site string) {

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

	if keys, err := vault.RevealEntry(site); err != nil {
		log.Fatal("Unable to add entry for: "+site, err)
	} else {
		for key, val := range keys {
			fmt.Printf("%s=%s\n", key, val)
		}
	}
}
