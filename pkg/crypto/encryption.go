package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

// EncryptionService handles encryption/decryption of sensitive data
type EncryptionService struct {
	key []byte
}

// NewEncryptionService creates a new encryption service with the given key
func NewEncryptionService(secretKey string) *EncryptionService {
	// Create a 32-byte key from the secret key using SHA-256
	hash := sha256.Sum256([]byte(secretKey))
	return &EncryptionService{
		key: hash[:],
	}
}

// Encrypt encrypts the given plaintext using AES-GCM
func (e *EncryptionService) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts the given ciphertext using AES-GCM
func (e *EncryptionService) Decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext_bytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext_bytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// EncryptMT5Credentials encrypts MT5 account credentials
func (e *EncryptionService) EncryptMT5Credentials(login, password, server string) (string, error) {
	credentials := fmt.Sprintf("%s:%s:%s", login, password, server)
	return e.Encrypt(credentials)
}

// DecryptMT5Credentials decrypts MT5 account credentials
func (e *EncryptionService) DecryptMT5Credentials(encryptedCredentials string) (login, password, server string, err error) {
	decrypted, err := e.Decrypt(encryptedCredentials)
	if err != nil {
		return "", "", "", err
	}

	// Split the credentials
	parts := strings.Split(decrypted, ":")
	if len(parts) != 3 {
		return "", "", "", errors.New("invalid credentials format")
	}

	return parts[0], parts[1], parts[2], nil
}
