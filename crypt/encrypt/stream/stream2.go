package stream

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
	"log"
	"os"
	"runtime/debug"
	"strings"

	"github.com/colt3k/utils/crypt"
	"github.com/colt3k/utils/crypt/encrypt/scrypt"
	"github.com/colt3k/utils/encode"
	"github.com/colt3k/utils/encode/encodeenum"
)

var (
	ScryptParams = scrypt.Params{N: 65536, R: 1, P: 2, SaltLen: 16, DKLen: 32}
	encSalt      string
	salt         []byte
)

func GenSalt() []byte {
	if len(encSalt) <= 0 {
		saltAR := crypt.GenSalt(nil, ScryptParams.SaltLen)
		ScryptParams.Salt = saltAR
		salt = saltAR
		// Save off SALT
		encodedSalt := encode.Encode(saltAR, encodeenum.B64STD)
		encSalt = encodedSalt
	}
	return salt
}

/*
EncToTemp encrypt file to a temp file
*/
func EncToTemp(filepath string, pass []byte) string {

	fo, err := os.Open(filepath)
	if err != nil {
		stackInfo := debug.Stack()
		log.Fatalf("err opening file, %v\n %v\n%v", filepath, err, string(stackInfo))
	}
	defer fo.Close()

	GenSalt()
	derivedKey, err := scrypt.Key(string(pass), ScryptParams)
	if err != nil {
		panic(err)
	}
	parts := strings.Split(string(derivedKey), "$")
	lastPart := parts[len(parts)-1]
	dk := []byte(lastPart)
	aesKey2 := dk[0:16]
	aesIv2 := dk[16:32]
	block, err := aes.NewCipher(aesKey2)
	if err != nil {
		panic(err)
	}

	stream := cipher.NewOFB(block, aesIv2)

	tf, err := os.CreateTemp(os.TempDir(), "ctcloud")
	if err != nil {
		stackInfo := debug.Stack()
		log.Fatalf("err opening temp file, %v\n %v\n%v", filepath, err, string(stackInfo))
	}

	// var out bytes.Buffer
	writer := &cipher.StreamWriter{S: stream, W: tf}
	// Copy the input to the output buffer, encrypting as we go.
	if _, err := io.Copy(writer, fo); err != nil {
		panic(err)
	}

	tf.Close()

	return tf.Name()
}

// DecFromTemp decrypt from temp file
// func DecFromTemp(tmpFile string, pass []byte, saveto string, salt string) {
// 	f := filenative.NewFile(tmpFile)
//
// 	fo, err := os.Open(f.Path())
// 	if err != nil {
// 		stackInfo := debug.Stack()
// 		log.Fatalf("err opening file, %v\n %v\n%v", tmpFile, err, string(stackInfo))
// 	}
// 	defer fo.Close()
//
// 	encSalt = salt
// 	decodeSalt := encode.Decode([]byte(encSalt), encodeenum.B64STD)
// 	ScryptParams.Salt = decodeSalt
//
// 	derivedKey, err := scrypt.Key(string(pass), ScryptParams)
// 	if err != nil {
// 		panic(err)
// 	}
// 	parts := strings.Split(string(derivedKey), "$")
// 	lastPart := parts[len(parts)-1]
// 	dk := []byte(lastPart)
// 	aesKey2 := dk[0:16]
// 	aesIv2 := dk[16:32]
// 	block, err := aes.NewCipher(aesKey2)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	stream := cipher.NewOFB(block, aesIv2)
//
// 	tf, err := os.Create(saveto)
// 	if err != nil {
// 		stackInfo := debug.Stack()
// 		log.Fatalf("err opening saveto file, %v\n%v", err, string(stackInfo))
// 	}
// 	defer tf.Close()
//
// 	reader := &cipher.StreamReader{S: stream, R: fo}
// 	// Copy the input to the output stream, decrypting as we go.
// 	if _, err := io.Copy(tf, reader); err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Written to file path decrypted")
// }
