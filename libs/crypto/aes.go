package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"log"

	"github.com/zenazn/pkcs7pad"
)

func EncryptAES(key []byte, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return handleEncryptionError(err)
	}

	pPlaintext := pkcs7pad.Pad(plaintext, block.BlockSize())

	ciphertext := make([]byte, aes.BlockSize+len(pPlaintext))
	iv := ciphertext[:aes.BlockSize]

	cipher.NewCBCEncrypter(block, iv).CryptBlocks(ciphertext[aes.BlockSize:], pPlaintext)

	return ciphertext, nil
}

func DecryptAES(key []byte, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return handleDecryptionError(err)
	}

	pPlaintext := make([]byte, len(ciphertext))
	iv := ciphertext[:aes.BlockSize]

	cipher.NewCBCDecrypter(block, iv).CryptBlocks(pPlaintext[aes.BlockSize:], ciphertext[aes.BlockSize:])

	plaintext, err := pkcs7pad.Unpad(pPlaintext[aes.BlockSize:])
	if err != nil {
		return handleDecryptionError(err)
	}

	return plaintext, nil
}

func handleEncryptionError(err error) ([]byte, error) {
	log.Printf("Error in EncryptionAES: %v\n", err)
	return nil, err
}

func handleDecryptionError(err error) ([]byte, error) {
	log.Printf("Error in DecryptionAES: %v\n", err)
	return nil, err
}
