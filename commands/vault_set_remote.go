package commands

import (
	"log"

	"github.com/jandre/passward/passward"
)

func VaultSetRemote(name string, url string) {

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

	if url == "" {
		log.Fatal("url is required.")
	}

	err = vault.SetRemote(url)
	if err != nil {
		log.Fatal("Unable to set vault remote: ", err)
	}

	err = vault.Sync()
	if err != nil {
		log.Fatal("Unable to set vault remote: ", err)
	}

}
