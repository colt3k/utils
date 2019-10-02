package crypt

import (
	"crypto/rand"
)

/*
GenSalt generate a salt if one doesn't exist of the length specified
*/
func GenSalt(salt []byte, length int) []byte {
	if salt == nil || len(salt) < length {
		//log.Logln(log.DEBUG, "Generating a Salt, none passed in.")
		b := GenerateRandomBytes(length)
		salt = b
	}
	return salt
}

/*
GenerateRandomBytes returns securely generated random bytes.
It will return an error if the system's secure random
number generator fails to function correctly, in which
case the caller should not continue.
*/
func GenerateRandomBytes(length int) []byte {
	data := make([]byte, length)
	_, err := rand.Read(data)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		panic("unable to read enough random bytes for iv")
	}

	return data
}
