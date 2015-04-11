package commands

import (
	"fmt"
	"log"

	"github.com/jandre/passward/passward"
)

func formatVaultListing(vaults *map[string]*passward.Vault) {

	if len(*vaults) == 0 {
		fmt.Println("No vaults found.  Add a new vault with `passward vault new <name>`.")
	} else {
		fmt.Printf("Found %d vaults:\n", len(*vaults))

		// todo: also print remote, etc
		for name, _ := range *vaults {
			fmt.Printf("\t%s\n", name)
		}
	}
}

func VaultList() {

	passwardPath := passward.DetectPasswardPath()
	pw, err := passward.ReadPassward(passwardPath)

	if err != nil {
		log.Fatal("There was a problem loading the configuration. Did you run `passward setup?`", err)
	}

	vaults := pw.GetVaults()

	formatVaultListing(&vaults)

}
