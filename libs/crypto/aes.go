package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"log"
)

func EncryptAES(key []byte, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Printf("Error in EncryptionAES: %v\n", err)
		return nil, err
	}

	pPlaintext := PKCS5Padding(plaintext, block.BlockSize(), len(plaintext))

	ciphertext := make([]byte, aes.BlockSize+len(pPlaintext))
	iv := ciphertext[:aes.BlockSize]
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(ciphertext[aes.BlockSize:], pPlaintext)

	return ciphertext, nil
}

func DecryptAES(key []byte, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Printf("Error in DecryptionAES: %v\n", err)
		return nil, err
	}
	plaintext := make([]byte, len(ciphertext))
	iv := ciphertext[:aes.BlockSize]

	cipher.NewCBCDecrypter(block, iv).CryptBlocks(plaintext, ciphertext[aes.BlockSize:])
	return plaintext, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int, after int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
