package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/jandre/passward/passward"
	"github.com/segmentio/go-prompt"
)

func VaultFetch(url string, name string) {

	passwardPath := passward.DetectPasswardPath()

	pw, err := passward.ReadPassward(passwardPath)

	if err != nil {
		log.Fatal("There was a problem loading the configuration. Did you run `passward setup?`", err)
	}

	passphrase := prompt.PasswordMasked("Enter your passphrase to unlock your keys (empty for none)")

	if err := pw.Unlock(passphrase); err != nil {
		log.Fatal("Invalid passphrase.", err)
	}

	vault, err := pw.FetchVault(url, name)

	if err != nil {
		fmt.Println("Unable fetch vault from remote:", url)
		fmt.Println("Error is:", err)
		fmt.Println("If authentication fails, you may also need to add the following ssh public key to the remote git server:")
		fmt.Println("")
		fmt.Println(pw.Credentials.PublicKeyString())
		os.Exit(1)
	}

	err = pw.UseVault(vault.Name)

	if err != nil {
		fmt.Println("Vault downloaded successfully, but unable to switch to vault.")
	}

	fmt.Printf("Vault fetched successfully: %s\n", vault.Name)
	fmt.Println("We have automatically switched to this as the active vault.  You can select another vault using `vault use`.")
}
