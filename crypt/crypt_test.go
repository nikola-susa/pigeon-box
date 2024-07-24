package crypt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncryptDecryptWithValidPassphrase(t *testing.T) {
	passphrase := "passphrase"
	plaintext := []byte("plaintext")

	ciphertext, err := Encrypt(passphrase, plaintext)
	assert.NoError(t, err)

	decrypted, err := Decrypt(passphrase, ciphertext)
	assert.NoError(t, err)

	assert.Equal(t, plaintext, decrypted)
}

func TestEncryptWithEmptyPassphrase(t *testing.T) {
	passphrase := ""
	plaintext := []byte("plaintext")

	_, err := Encrypt(passphrase, plaintext)
	assert.Error(t, err)
}

func TestDecryptWithInvalidPassphrase(t *testing.T) {
	passphrase := "passphrase"
	invalidPassphrase := "passphrase2"
	plaintext := []byte("plaintext")

	ciphertext, err := Encrypt(passphrase, plaintext)
	assert.NoError(t, err)

	_, err = Decrypt(invalidPassphrase, ciphertext)
	assert.Error(t, err)
}

func TestDecryptWithModifiedCiphertext(t *testing.T) {
	passphrase := "passphrase"
	plaintext := []byte("plaintext")

	ciphertext, err := Encrypt(passphrase, plaintext)
	assert.NoError(t, err)

	ciphertext[0] ^= 0xff

	_, err = Decrypt(passphrase, ciphertext)
	assert.Error(t, err)
}
