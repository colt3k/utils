package aesvault

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

//AES store variables for encryption
type AES struct {
	AESKeyLength   int
	SaltLength     int
	IVLength       int
	IterationCount int
	Plaintext      string
	CipherText     string
	CipherKey      []byte
	HMACKey        []byte
	IV             []byte
}

func (a *AES) Validate() bool {
	return false
}

//Encrypt interface method for encryption
func (a *AES) Encrypt() ([]byte, error) {
	aesCipher, err := aes.NewCipher(a.CipherKey)
	if err != nil {
		return nil, err
	}

	//Add Padding
	//content := padding.PKCS5Padding([]byte(a.plaintext), aesCipher.BlockSize())
	fmt.Printf("Encrypting:|%s|\n", a.Plaintext)
	content := pad([]byte(a.Plaintext))
	fmt.Printf("Padded: |%s|\n",content)
	ciphertext := make([]byte, len(content))

	aesBlock := cipher.NewCTR(aesCipher, a.IV)
	aesBlock.XORKeyStream(ciphertext, content)
	fmt.Printf("Cipher: |%v|\n",ciphertext)
	return ciphertext, nil
}

//Decrypt interface method for decryption
func (a *AES) Decrypt() (string, error) {
	aesCipher, err := aes.NewCipher(a.CipherKey)
	if err != nil {
		return "", err
	}

	plainText := make([]byte, len(a.CipherText))

	aesBlock := cipher.NewCTR(aesCipher, a.IV)
	aesBlock.XORKeyStream(plainText, []byte(a.CipherText))

	result, err := unpad(plainText)
	if err != nil {
		panic(err)
	}

	return string(result), nil
}

func pad(src []byte) []byte {
	padlen := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padlen)}, padlen)
	return append(src, padtext...)
}
func unpad(src []byte) ([]byte, error) {
	length := len(src)
	padlen := int(src[length-1])
	if padlen > length {
		return nil, fmt.Errorf("invalid padding")
	}
	return src[:(length - padlen)], nil
}
