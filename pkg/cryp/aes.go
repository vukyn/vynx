package cryp

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
)

// AES-GCM encryption
func EncryptAES(plaintext string, aesKey string) (string, error) {
	if aesKey == "" {
		return "", errors.New("aesKey is empty")
	}
	block, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)
	// Store nonce + ciphertext
	result := append(nonce, ciphertext...)
	return base64.StdEncoding.EncodeToString(result), nil
}

// AES-GCM decryption
func DecryptAES(ciphertextB64 string, aesKey string) (string, error) {
	if aesKey == "" {
		return "", errors.New("aesKey is empty")
	}
	data, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(data) < gcm.NonceSize() {
		return "", err
	}
	nonce := data[:gcm.NonceSize()]
	ciphertext := data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
