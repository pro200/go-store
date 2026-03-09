package lib

import (
	"github.com/pro200/go-aes256"
)

var (
	encryptKey, _ = CUID()
	encryptIv     = encryptKey[3:19]
)

func Encrypt(data []byte) ([]byte, error) {
	return aes256.Encrypt([]byte(encryptKey), []byte(encryptIv), data)
}

func Decrypt(encrypted []byte) ([]byte, error) {
	return aes256.Decrypt([]byte(encryptKey), []byte(encryptIv), encrypted)
}
