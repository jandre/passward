package commands

import (
	"log"

	"github.com/jandre/passward/passward"
	"github.com/segmentio/go-prompt"
)

func VaultSecretAdd(name string, key string, value string) {

	passwardPath := passward.DetectPasswardPath()

	pw, err := passward.ReadPassward(passwardPath)

	if err != nil {
		log.Fatal("There was a problem loading the configuration. Did you run `passward setup?`", err)
	}

	vault := pw.GetVault(name)

	if vault == nil {
		log.Fatal("Vault not found: " + name)
	}

	passphrase := prompt.PasswordMasked("Enter your passphrase to unlock your keys (empty for none)")
	pw.Unlock(passphrase)

}
