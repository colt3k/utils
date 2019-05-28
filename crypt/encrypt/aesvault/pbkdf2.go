package aesvault

import (
	"crypto/sha256"

	"github.com/colt3k/utils/hash/h_mac"
	"golang.org/x/crypto/pbkdf2"
)

func PBKDF2KeySHA256(password *string, salt []byte, a *AES) *h_mac.Key {

	k := pbkdf2.Key([]byte(*password), salt, a.IterationCount, 2*a.AESKeyLength+a.IVLength, sha256.New)

	return &h_mac.Key{
		CipherKey: k[:a.AESKeyLength],
		HMACKey:   k[a.AESKeyLength:(a.AESKeyLength * 2)],
		IV:        k[(a.AESKeyLength * 2) : (a.AESKeyLength*2)+a.IVLength],
	}
}
