package commands

import (
	"log"

	"github.com/jandre/passward/passward"
)

func chooseVault(pw *passward.Passward, name string) *passward.Vault {
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

	return vault
}
