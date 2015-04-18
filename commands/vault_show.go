package commands

import (
	"fmt"
	"log"

	"github.com/jandre/passward/passward"
)

func VaultShow(name string) {
	var vault *passward.Vault

	passwardPath := passward.DetectPasswardPath()

	pw, err := passward.ReadPassward(passwardPath)

	if err != nil {
		log.Fatal("There was a problem loading the configuration. Did you run `passward setup?`", err)
	}

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

	users := vault.Users()
	fmt.Printf("Showing vault: %s\n", vault.Name)
	fmt.Printf("-- Found %d users\n", len(users))

	for _, user := range users {
		fmt.Printf("\tUser: %s\n", user.Email())
	}

	entries := vault.Entries()
	fmt.Printf("-- Found %d sites\n", len(entries))

	for _, user := range entries {
		fmt.Printf("\tSite: %s\n", user.Name())
	}
}
