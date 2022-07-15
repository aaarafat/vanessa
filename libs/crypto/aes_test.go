package crypto

import (
	"crypto/aes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createKey() []byte {
	key := []byte("WnZr4u7x!A%D*G-K")
	return key
}

func TestEncryptAES(t *testing.T) {
	key := createKey()
	testCases := []struct {
		plaintext string
	}{
		{"Hello, world!"},
		{"Hello"},
		{"He"},
		{"h"},
		{"heeeeeeeeeeeeeeeeeeeeeeeee, worlddddddddddddddddd!"},
		{string([]byte{0, 4})},
		{string([]byte{4, 0})},
	}

	for _, testCase := range testCases {
		plaintext := []byte(testCase.plaintext)
		ciphertext, err := EncryptAES(key, plaintext)
		assert.NoError(t, err)
		assert.NotEqual(t, plaintext, ciphertext)
		assert.Equal(t, len(ciphertext)%aes.BlockSize, 0)
	}
}

func TestDecryptAES(t *testing.T) {
	key := createKey()

	testCases := []struct {
		plaintext string
	}{
		{"Hello, world!"},
		{"Hello"},
		{"He"},
		{"h"},
		{"heeeeeeeeeeeeeeeeeeeeeeeee, worlddddddddddddddddd!"},
		{string([]byte{0, 4})},
		{string([]byte{4, 0})},
	}

	for _, testCase := range testCases {
		ciphertext, err := EncryptAES(key, []byte(testCase.plaintext))
		assert.Nil(t, err)
		decrypted, err := DecryptAES(key, ciphertext)
		assert.Nil(t, err)
		assert.Equal(t, []byte(testCase.plaintext), decrypted)
	}
}
