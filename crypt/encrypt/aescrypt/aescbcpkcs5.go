package aescrypt

/*
ENC_ALGORITHM		: AES
ENC_TRANSFORMATION	: AES/CBC/PKCS5Padding	(Cipher Instance)
*/
import (
	"crypto/aes"
	"crypto/cipher"
	"time"

	log "github.com/colt3k/nglog/ng"
	
	"github.com/colt3k/utils/crypt/encrypt/padding"
)

var (
	aesKeyLength = 32               // double of the SALT
	saltLength   = aesKeyLength / 2 // half of the key length
	ivSize       = aes.BlockSize

	iterationCount = 65536
)

//AES store variables for encryption
type AES struct {
	aesKeyLength   int
	saltLength     int
	ivSize         int
	iterationCount int
	plaintext      string
	cipherText     []byte
	key            []byte
}

//New create a new instance of the AES struct
func New(keylength, saltLength, ivSize, iterations int, plaintext *string, cipherText *[]byte, key *[]byte) *AES {
	a := &AES{}

	if keylength > -1 {
		a.aesKeyLength = keylength
	}
	if saltLength > -1 {
		a.saltLength = saltLength
	}
	if ivSize > -1 {
		a.ivSize = ivSize
	}
	if iterations > -1 {
		a.iterationCount = iterations
	}
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

	start := time.Now()
	// Now we need to extract AES key and IV from newly derived key
	aesKey2 := derivedKey[0:16]
	aesIv2 := derivedKey[16:32]

	block, err := aes.NewCipher(aesKey2)
	if err != nil {
		panic(err)
	}

	/** DOING IT WITHOUT OPENSSL REQUIRES BLOCKS OF 32, unless using PKCS5 Padding(also referred to PKCS7) ***/
	encMode := cipher.NewCBCEncrypter(block, aesIv2)

	//Convert string to byte[]
	content := []byte(a.plaintext)
	//Add Padding
	content = padding.PKCS5Padding(content, block.BlockSize())
	//Create byte[] the size of the content
	crypted := make([]byte, len(content))

	encMode.CryptBlocks(crypted, content)
	elapsed := time.Since(start)
	log.Logf(log.DEBUG, "Execution took %s", elapsed)

	return crypted

}

func aesDecrypt(a *AES, derivedKey []byte) []byte {
	start := time.Now()
	// Now we need to extract AES key and IV from newly derived key
	aesKey2 := derivedKey[0:16]
	aesIv2 := derivedKey[16:32]

	block, err := aes.NewCipher(aesKey2)
	if err != nil {
		panic(err)
	}

	/** DOING IT WITHOUT OPENSSL REQUIRES BLOCKS OF 32, unless using PKCS5 Padding ***/
	decMode := cipher.NewCBCDecrypter(block, aesIv2)

	decrypted := make([]byte, len(a.cipherText))

	decMode.CryptBlocks(decrypted, a.cipherText)

	elapsed := time.Since(start)
	log.Logf(log.DEBUG, "Execution took %s", elapsed)
	padded := padding.PKCS5Trimming(decrypted)
	return padded
}
