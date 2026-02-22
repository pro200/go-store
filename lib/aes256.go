package lib

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

var (
	encryptKey, _ = CUID()
	encryptIv     = encryptKey[3:19]
)

// PKCS7 Padding
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// PKCS7 Unpadding
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, errors.New("invalid padding size")
	}

	padding := int(data[len(data)-1])
	if padding == 0 || padding > blockSize {
		return nil, errors.New("invalid padding")
	}

	return data[:len(data)-padding], nil
}

// AES-CBC Encrypt
func AesEncrypt(key, iv, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(iv) != aes.BlockSize {
		return nil, errors.New("invalid IV size")
	}

	data = pkcs7Pad(data, aes.BlockSize)

	encrypted := make([]byte, len(data))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encrypted, data)

	return encrypted, nil
}

// AES-CBC Decrypt
func AesDecrypt(key, iv, encryted []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(iv) != aes.BlockSize {
		return nil, errors.New("invalid IV size")
	}

	if len(encryted)%aes.BlockSize != 0 {
		return nil, errors.New("encryted is not a multiple of block size")
	}

	data := make([]byte, len(encryted))

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(data, encryted)

	return pkcs7Unpad(data, aes.BlockSize)
}

func Encrypt(data []byte) ([]byte, error) {
	return AesEncrypt([]byte(encryptKey), []byte(encryptIv), data)
}

func Decrypt(encrypted []byte) ([]byte, error) {
	return AesDecrypt([]byte(encryptKey), []byte(encryptIv), encrypted)
}

//func main() {
//	key := []byte("1234567890123456") // 16, 24, 32 bytes
//	iv := []byte("abcdefghijklmnop")  // 16 bytes
//
//	data := []byte("hello world")
//
//	enc, _ := AesEncrypt(key, iv, data)
//	dec, _ := AesDecrypt(key, iv, enc)
//
//	fmt.Println(string(dec)) // hello world
//}
