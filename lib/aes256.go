package lib

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

var (
	encryptKey, _ = CUID()
	encryptIv     = encryptKey[3:19]
)

// AesEncrypt encrypts given text in AES 256 CBC
//
//	key := "ovmktblgkdcsvtibhqcwmegxpjtmqlwh" 23byte
//	iv := "kh5b2kldgfo3ng4k"
//	text := "Hello, World!"
//	encrypted, _ := AesEncrypt(key, iv, text)
//	fmt.Println(encrypted) -> Wn1j1g5CIzxOetrdehyrLQ==
func AesEncrypt(key, iv, plaintext string) (string, error) {
	var plainTextBlock []byte
	length := len(plaintext)

	if length%16 != 0 {
		extendBlock := 16 - (length % 16)
		plainTextBlock = make([]byte, length+extendBlock)
		copy(plainTextBlock[length:], bytes.Repeat([]byte{uint8(extendBlock)}, extendBlock))
	} else {
		plainTextBlock = make([]byte, length)
	}

	copy(plainTextBlock, plaintext)

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, len(plainTextBlock))
	mode := cipher.NewCBCEncrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, plainTextBlock)
	str := base64.StdEncoding.EncodeToString(ciphertext)

	return str, nil
}

func Encrypt(plaintext string) (string, error) {
	return AesEncrypt(encryptKey, encryptIv, plaintext)
}

// AesDecrypt decrypts given text in AES 256 CBC
//
//	key := "ovmktblgkdcsvtibhqcwmegxpjtmqlwh" 32byte
//	iv := "kh5b2kldgfo3ng4k"
//	encrypted := "Wn1j1g5CIzxOetrdehyrLQ=="
//	decrypted, _ := AesDecrypt(key, iv, encrypted)
//	fmt.Println(decrypted) -> Hello, World!
func AesDecrypt(key, iv, encrypted string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return "", fmt.Errorf("block size cant be zero")
	}

	mode := cipher.NewCBCDecrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, ciphertext)

	// PKCS5UnPadding pads a certain blob of data with necessary data to be used in AES block cipher
	length := len(ciphertext)
	unpadding := int(ciphertext[length-1])
	ciphertext = ciphertext[:(length - unpadding)]

	return string(ciphertext), nil
}

func Decrypt(encrypted string) (string, error) {
	return AesDecrypt(encryptKey, encryptIv, encrypted)
}
