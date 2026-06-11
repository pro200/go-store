package store

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/pro200/go-aes256"
)

// deriveKey converts a machine identifier into a 32-byte AES-256 key.
// The domain-separation prefix prevents key reuse with other tools
// that hash the same machine id.
func deriveKey(machineID string) []byte {
	sum := sha256.Sum256([]byte("github.com/pro200/go-store:v2:" + machineID))
	return sum[:]
}

// encrypt returns a random 16-byte IV followed by the AES-256-CBC
// ciphertext of plaintext.
func encrypt(key, plaintext []byte) ([]byte, error) {
	iv := make([]byte, aes256.IVSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, fmt.Errorf("generate iv: %w", err)
	}

	ciphertext, err := aes256.Encrypt(key, iv, plaintext)
	if err != nil {
		return nil, err
	}
	return append(iv, ciphertext...), nil
}

// decrypt reverses encrypt: it splits the leading IV from data and
// decrypts the remainder.
func decrypt(key, data []byte) ([]byte, error) {
	if len(data) < aes256.IVSize {
		return nil, errors.New("encrypted data too short")
	}
	return aes256.Decrypt(key, data[:aes256.IVSize], data[aes256.IVSize:])
}
