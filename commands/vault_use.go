package commands

import (
	"fmt"
	"log"

	"github.com/jandre/passward/passward"
)

func VaultUse(name string) {

	passwardPath := passward.DetectPasswardPath()

	pw, err := passward.ReadPassward(passwardPath)

	if err != nil {
		log.Fatal("There was a problem loading the configuration. Did you run `passward setup?`", err)
	}

	err = pw.UseVault(name)
	if err != nil {
		fmt.Println("No vault found: " + name)
	} else {
		if err := pw.Save(); err != nil {
			log.Fatal("Unable to select vault:", err)
		} else {
			fmt.Println("Vault selected: " + name)
		}
	}

}
