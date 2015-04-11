package commands

import (
	"fmt"
	"log"

	"github.com/jandre/passward/passward"
)

func VaultSecretAdd(name string, key string, value string) {

	passwardPath := passward.DetectPasswardPath()

	pw, err := passward.ReadPassward(passwardPath)

	if err != nil {
		log.Fatal("There was a problem loading the configuration. Did you run `passward setup?`", err)
	}

	if err = pw.AddVault(name); err != nil {
		log.Fatal("Error creating vault: ", err)
	}

	fmt.Println("Successfully created new vault: ", name)

}
