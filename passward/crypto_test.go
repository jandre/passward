package passward

import (
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {

	passphrase := "my secret passphrase"
	secret := "secret"

	ciphertext, err := EncryptAndBase64String(passphrase, secret)

	if err != nil {
		t.Fatal(err)
	}

	cleartext, err := DecryptBase64String(passphrase, ciphertext)

	if err != nil {
		t.Fatal(err)
	}

	if cleartext != secret {
		t.Fatal("mismatch:", cleartext, secret)
	}
}
