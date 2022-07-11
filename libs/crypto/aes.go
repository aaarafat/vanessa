package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"log"
)

func EncryptAES(key []byte, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return handleEncryptionError(err)
	}

	pPlaintext, err := PKCS7Padding(plaintext, block.BlockSize())
	if err != nil {
		return handleEncryptionError(err)
	}

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

	pPlaintext := make([]byte, len(ciphertext)-aes.BlockSize)
	iv := ciphertext[:aes.BlockSize]

	cipher.NewCBCDecrypter(block, iv).CryptBlocks(pPlaintext, ciphertext[aes.BlockSize:])

	plaintext, err := PKCS7UnPadding(pPlaintext, block.BlockSize())
	if err != nil {
		return handleDecryptionError(err)
	}

	return plaintext, nil
}

// https://gist.github.com/nanmu42/b838acc10d393bc51cb861128ce7f89c
func PKCS7UnPadding(data []byte, blockSize int) ([]byte, error) {
	padding := int(data[len(data)-1])
	ref := bytes.Repeat([]byte{byte(padding)}, padding)
	if padding > blockSize || padding == 0 || !bytes.HasSuffix(data, ref) {
		return nil, fmt.Errorf("PKCS7: Invalid padding, padding %d, blocksize %d\n", padding, blockSize)
	}
	return data[:len(data)-padding], nil
}

func PKCS7Padding(data []byte, blockSize int) ([]byte, error) {
	if blockSize < 0 || blockSize > 256 {
		return nil, fmt.Errorf("PKCS7: Invalid blocksize %d\n", blockSize)
	}
	padding := (blockSize - len(data)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...), nil
}

func handleEncryptionError(err error) ([]byte, error) {
	log.Printf("Error in EncryptionAES: %v\n", err)
	return nil, err
}

func handleDecryptionError(err error) ([]byte, error) {
	log.Printf("Error in DecryptionAES: %v\n", err)
	return nil, err
}
