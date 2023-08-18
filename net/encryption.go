package net

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

const DecryptError = "decrypt Failed, invalid key"

func (net Net) encrypt(text string) (string, error) {
	if net.EncryptionKey == nil {
		return text, nil
	}
	block, err := aes.NewCipher(net.EncryptionKey[:])
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	return string(gcm.Seal(nonce, nonce, []byte(text), nil)), nil

}

func (net Net) decrypt(text string) (string, error) {
	if net.EncryptionKey == nil {
		return text, nil
	}
	ciphertext := []byte(text)
	block, err := aes.NewCipher(net.EncryptionKey[:])
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(ciphertext) < gcm.NonceSize() {
		return "", errors.New("malformed ciphertext")
	}
	data, err := gcm.Open(nil, ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():], nil)
	return string(data), err
}

// Pad the input to a multiple of the block size using PKCS7 padding
func pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// Remove PKCS7 padding from the input
func unpad(data []byte) ([]byte, error) {
	padding := int(data[len(data)-1])
	if padding > len(data) {
		return nil, errors.New(DecryptError)
	}
	return data[:len(data)-padding], nil
}

// Pad the key to 32 bytes
func padKey(key []byte) []byte {
	paddedKey := make([]byte, 32)
	copy(paddedKey, key)
	return paddedKey
}

func NewKey(key string) []byte {
	if key == "" {
		return nil
	}
	return padKey([]byte(key))
}
