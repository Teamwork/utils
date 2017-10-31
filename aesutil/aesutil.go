// Package aesutil provides a set of functions for 2-way encryption and hashing.
package aesutil // import "github.com/teamwork/utils/aesutil"

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// Encrypt binary data to a base64 string with AES using the key provided.
func Encrypt(keyString string, data []byte) (string, error) {
	key := []byte(keyString)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt a base64 string to binary data with AES using the key provided.
func Decrypt(keyString string, base64Data string) ([]byte, error) {
	key := []byte(keyString)
	ciphertext, err := base64.URLEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("Ciphertext provided is smaller than AES block size")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}
