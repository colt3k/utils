package aescrypt_gcm

/*
GCM (Galois/Counter Mode): designed to provide both data authenticity (integrity) and confidentiality
	defined for block ciphers with a block size of 128 bits.
	encrypted text then contains the IV, ciphertext, and authentication tag

Implementation: GCM with IV/Nonce of 16 and Padding
	Padding is added/removed in order to hide the size of original text
	Standard Tag size 16 -> 128 bit
	Standard IV/Nonce is 12, changed to 16

Test
	Demos use of GCM with SCrypt
	Note Scrypt Parameters are set explicitly so key doesn't change from system to system and tests are valid
*/
import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"github.com/colt3k/utils/encode"
	"github.com/colt3k/utils/encode/encodeenum"
)

var (
	aesKeyLength = 32               // double of the SALT
)

//AES store variables for encryption
type AES struct {
	plaintext      string
	cipherText     []byte
	key            []byte
}

func New(plaintext *string, cipherText *[]byte, key *[]byte) *AES {
	a := &AES{}

	if plaintext != nil && len(*plaintext) > 0 {
		a.plaintext = *plaintext
	}
	if cipherText != nil {
		a.cipherText = *cipherText
	}
	if key != nil {
		a.key = *key
	}

	return a
}

func (a *AES) Validate() bool {
	return false
}

//Encrypt interface method for encryption
func (a *AES) Encrypt() []byte {
	return aesCrypt(a, a.key)
}

//Decrypt interface method for decryption
func (a *AES) Decrypt() []byte {
	return aesDecrypt(a, a.key)
}

func aesCrypt(a *AES, derivedKey []byte) []byte {

	// Now we need to extract AES key and IV from newly derived key
	aesKey2 := derivedKey[0:16]
	aesIv2 := derivedKey[16:32]

	block, err := aes.NewCipher(aesKey2)
	if err != nil {
		panic(err.Error())
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	//iv := make([]byte, 12)
	//if _, err := io.ReadFull(rand.Reader, iv); err != nil {
	//	panic(err.Error())
	//}

	// Tag Size 16 is 128 bit, this is default for NewGCM() also
	aesgcm, err := cipher.NewGCMWithNonceSize(block, 16)
	if err != nil {
		panic(err.Error())
	}

	//Convert string to byte[]
	content := []byte(a.plaintext)
	// ... NO PADDING NEEDED IN GCM but length of original Plaintext is exposed without it
	content = addPadding(content)

	ciphertext := aesgcm.Seal(nil, aesIv2, content, nil)

	return ciphertext
}

func aesDecrypt(a *AES, derivedKey []byte) []byte {

	// Now we need to extract AES key and IV from newly derived key
	aesKey2 := derivedKey[0:16]
	aesIv2 := derivedKey[16:32]

	block, err := aes.NewCipher(aesKey2)
	if err != nil {
		panic(err.Error())
	}

	// Tag Size 16 is 128 bit, this is default for NewGCM() also
	aesgcm, err := cipher.NewGCMWithNonceSize(block, 16)
	if err != nil {
		panic(err.Error())
	}
	cipherText := a.cipherText
	//cipherText = removePadding(cipherText)

	plaintext, err := aesgcm.Open(nil, aesIv2, cipherText, nil)
	if err != nil {
		panic(err.Error())
	}
	plaintext = removePadding(plaintext)

	return plaintext
}

func addPadding(msg []byte) []byte {
	AES_BLOCK_SIZE_BYTES := 16
	padSize := floorMod(-len(msg), AES_BLOCK_SIZE_BYTES)
	if padSize == 0 {
		padSize = AES_BLOCK_SIZE_BYTES
	}
	cpadding := make([]rune, padSize)
	for i := range cpadding {
		cpadding[i]=rune(padSize)
	}
	padding := []byte(string(cpadding))
	fmt.Printf("- Padding Size: %v\n", padSize)
	fmt.Printf("- Padding in hex: %v\n",encode.Encode(padding, encodeenum.Hex))
	b := bytes.NewBuffer(msg)
	b.Write(padding)

	return b.Bytes()
}

func removePadding(pmsg []byte) []byte {
	AES_BLOCK_SIZE_BYTES := 16
	var msg []byte

	valid := true
	if len(pmsg) % AES_BLOCK_SIZE_BYTES != 0 {
		valid = false
	}
	padsize := int(pmsg[len(pmsg)-1])
	if padsize > len(pmsg) {
		valid = false
	}
	fmt.Printf("- Padding Size: %v\n", padsize)

	for i := len(pmsg) - 1 ; i > len(pmsg) - padsize; i = i - 1 {
		if int(pmsg[i]) != padsize {
			valid = false
		}
	}

	if valid {
		msg = pmsg[:len(pmsg) - padsize]
		fmt.Printf("-text without padding in hex: %v\n", encode.Encode(msg, encodeenum.Hex))
	}
	return msg
}

func floorMod(a, b int)	 int {
	return (a % b + b) %b
}