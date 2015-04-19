package commands

import (
	"fmt"
	"log"

	"github.com/jandre/passward/passward"
	"github.com/segmentio/go-prompt"
)

func VaultSetRemote(name string, url string) {

	passwardPath := passward.DetectPasswardPath()

	pw, err := passward.ReadPassward(passwardPath)

	if err != nil {
		log.Fatal("There was a problem loading the configuration. Did you run `passward setup?`", err)
	}

	if url == "" {
		log.Fatal("url is required.")
	}

	vault := chooseVault(pw, name)

	passphrase := prompt.PasswordMasked("Enter your passphrase to unlock your keys (empty for none)")
	if err := pw.Unlock(passphrase); err != nil {
		log.Fatal("Invalid passphrase.", err)
	}

	err = vault.SetRemote(url)
	if err != nil {
		log.Fatal("Unable to set vault remote: ", err)
	}

	fmt.Printf("Remote for `%s` successfully set to: %s\n", vault.Name, url)
	fmt.Printf("Will attempt to run initial sync...")
	err = vault.Sync()
	if err != nil {
		log.Fatal("Unable to sync vault remote: ", err)
	}
	fmt.Println("Vault sync'd successfully!")
}
