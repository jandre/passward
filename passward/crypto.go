//
// Encryption wrappers
//
package passward

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

func SignData(passphrase string, data string) string {
	key := []byte(passphrase)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func VerifySignature(passphrase string, data string, expected string) bool {
	signature := SignData(passphrase, data)
	b1, _ := base64.StdEncoding.DecodeString(signature)
	b2, _ := base64.StdEncoding.DecodeString(expected)
	return hmac.Equal(b1, b2)
}

//
// KeyGen() will generate a key with a passphrase that is `keySize` bytes
// in length.
//
func KeyGen(passphrase string, keySize uint) ([]byte, error) {

	// there's too many bytes requested
	// TODO: can generate multi-hashes
	if keySize > sha256.Size {
		return nil, errors.New("key size is too large: " + string(keySize))
	}

	// it's an n byte key, so let's generate a cryptographic hash of the passphrase and
	// use the first n bytes
	bytes := sha256.Sum256([]byte(passphrase))
	return bytes[:keySize], nil
}

//
// GenRandomIv() will generate random IV of `blockSize` bytes.
//
func GenRandomIv(blockSize int) ([]byte, error) {
	b := make([]byte, blockSize)
	bytesRead, err := rand.Read(b)

	if err != nil {
		return nil, err
	}

	if bytesRead != blockSize {
		return nil, errors.New("unable to generate random iv bytes")
	}

	return b, nil
}

//
// Encrypts a block with the given passphrase.
//
// Returns a byte block with <32 byte sha-256 hmac>< blockSize iv>< n blocks payload>
func Encrypt(passphrase string, bytes []byte) ([]byte, error) {

	key, err := KeyGen(passphrase, aes.BlockSize)

	if err != nil {
		return nil, err
	}

	aes, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	mode, err := cipher.NewGCM(aes)

	if err != nil {
		return nil, err
	}

	iv, err := GenRandomIv(mode.NonceSize())

	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	encrypted := mode.Seal(nil, iv, bytes, nil)

	result := append(iv, encrypted...)

	return result, nil
}

func Decrypt(passphrase string, data []byte) ([]byte, error) {

	key, err := KeyGen(passphrase, aes.BlockSize)

	if err != nil {
		return nil, err
	}

	aes, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	mode, err := cipher.NewGCM(aes)
	ivBytesSize := mode.NonceSize()
	iv := data[:ivBytesSize]
	cipherText := data[ivBytesSize:]

	if err != nil {
		return nil, err
	}

	return mode.Open(nil, iv, cipherText, nil)

}

//
// EncryptString encrypts a string `data` with the passphrase `passphrase`
//
func EncryptString(passphrase string, data string) ([]byte, error) {
	return Encrypt(passphrase, []byte(data))
}

//
// DecryptString
//
func DecryptString(passphrase string, data []byte) (string, error) {
	bytes, err := Decrypt(passphrase, data)

	if err != nil {
		return "", err
	}
	return string(bytes), err
}

func DecryptBase64String(passphrase string, input string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(input)

	if err != nil {
		return "", err
	}
	return DecryptString(passphrase, bytes)
}

//
// Encrypts and base64 encrypted output.
//
func EncryptAndBase64String(passphrase string, data string) (string, error) {
	result, err := EncryptString(passphrase, data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), nil
}
