package crypto

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createKey() []byte {
	key := make([]byte, 16)
	rand.Read(key)
	return key
}

func TestEncryptAES(t *testing.T) {
	key := createKey()

	plaintext := []byte("Hello, world!")
	ciphertext, err := EncryptAES(key, plaintext)

	assert.Nil(t, err)
	assert.NotEqual(t, plaintext, ciphertext)
}

func TestDecryptAES(t *testing.T) {
	key := createKey()

	plaintext := []byte("Hello, world!")
	ciphertext, _ := EncryptAES(key, plaintext)

	decrypted, err := DecryptAES(key, ciphertext)

	assert.Nil(t, err)
	assert.Equal(t, plaintext, decrypted)
}
