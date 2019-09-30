package stream

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/colt3k/nglog/ers/bserr"
	"github.com/colt3k/utils/encode"
	"github.com/colt3k/utils/encode/encodeenum"
	"github.com/colt3k/utils/file/filenative"

	"github.com/colt3k/utils/crypt"
	"github.com/colt3k/utils/crypt/encrypt/scrypt"
)

var (
	ScryptParams = scrypt.Params{N:65536, R:1, P:2, SaltLen:16, DKLen:32}
	encSalt string
	salt []byte
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
	bserr.StopErr(err, "err opening file, "+filepath)
	defer fo.Close()

	GenSalt()
	derivedKey, err := scrypt.Key(string(pass), ScryptParams)
	if err != nil {
		panic(err)
	}
	aesKey2 := derivedKey[0:16]
	aesIv2 := derivedKey[16:32]
	block, err := aes.NewCipher(aesKey2)
	if err != nil {
		panic(err)
	}

	stream := cipher.NewOFB(block, aesIv2)

	tf, err := ioutil.TempFile(os.TempDir(), "ctcloud")
	bserr.StopErr(err, "err opening temp file")

	//var out bytes.Buffer
	writer := &cipher.StreamWriter{S: stream, W: tf}
	// Copy the input to the output buffer, encrypting as we go.
	if _, err := io.Copy(writer, fo); err != nil {
		panic(err)
	}

	tf.Close()

	return tf.Name()
}

// DecFromTemp decrypt from temp file
func DecFromTemp(tmpFile string, pass []byte, saveto string, salt string) {
	f := filenative.NewFile(tmpFile)

	fo, err := os.Open(f.Path())
	bserr.StopErr(err, "err opening file, "+tmpFile)
	defer fo.Close()

	encSalt = salt
	decodeSalt := encode.Decode([]byte(encSalt), encodeenum.B64STD)
	ScryptParams.Salt = decodeSalt

	derivedKey, err := scrypt.Key(string(pass), ScryptParams)
	if err != nil {
		panic(err)
	}
	aesKey2 := derivedKey[0:16]
	aesIv2 := derivedKey[16:32]
	block, err := aes.NewCipher(aesKey2)
	if err != nil {
		panic(err)
	}

	stream := cipher.NewOFB(block, aesIv2)

	tf, err := os.Create(saveto)
	bserr.StopErr(err, "err opening saveto file")
	defer tf.Close()

	reader := &cipher.StreamReader{S: stream, R: fo}
	// Copy the input to the output stream, decrypting as we go.
	if _, err := io.Copy(tf, reader); err != nil {
		panic(err)
	}
	fmt.Println("Written to file path decrypted")
}