package commands

import (
	"fmt"
	"log"

	"github.com/jandre/passward/passward"
)

func VaultShow(name string) {

	passwardPath := passward.DetectPasswardPath()

	pw, err := passward.ReadPassward(passwardPath)

	if err != nil {
		log.Fatal("There was a problem loading the configuration. Did you run `passward setup?`", err)
	}

	// passphrase := prompt.PasswordMasked("Enter your passphrase to unlock your keys (empty for none)")

	// if err := pw.Unlock(passphrase); err != nil {
	// log.Fatal("Invalid passphrase.", err)
	// }

	vault := pw.GetVault(name)
	if vault == nil {
		log.Fatal("No vault found: ", name)
	}

	users := vault.Users()
	fmt.Printf("-- Found %d users\n", len(users))

	for _, user := range users {
		fmt.Println("\tUser:", user.Email())
	}
}
