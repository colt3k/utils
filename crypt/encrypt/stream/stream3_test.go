package stream

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/colt3k/utils/file/filenative"

	"github.com/colt3k/utils/crypt"
	"github.com/colt3k/utils/crypt/encrypt/scrypt"
)

func TestEncrypt(t *testing.T) {

	testfile := "/Users/gcollins/Desktop/b2test/hello1.txt"
	fo, err := os.Open(testfile)
	if err != nil {
		panic(err)
	}

	tf, err := ioutil.TempFile(os.TempDir(), "ctcloud")
	if err != nil {
		panic(err)
	}

	saltAR := crypt.GenSalt(nil, ScryptParams.SaltLen)
	ScryptParams.Salt = saltAR
	derivedKey, err := scrypt.Key(string("mypass"), ScryptParams)
	if err != nil {
		panic(err)
	}

	aesKey2 := derivedKey[0:16]
	//aesIv2 := derivedKey[16:32]

	hmacKey := []byte("this is my hmackey")
	err = Encrypt(fo, tf, aesKey2, hmacKey)
	if err != nil {
		panic(err)
	}

	fmt.Println("Temp at ", tf.Name())
	fo.Close()
	tf.Close()

	// DECRYPT
	f := filenative.NewFile(tf.Name())

	fo, err = os.Open(f.Path())
	if err != nil {
		panic(err)
	}

	tf, err = os.Create("./hello1.txt")
	if err != nil {
		panic(err)
	}

	err = Decrypt(fo, tf, aesKey2, hmacKey)
	if err != nil {
		panic(err)
	}
}

//type devZero byte
//
//func (z devZero) Read(b []byte) (int, error) {
//	for i := range b {
//		b[i] = byte(z)
//	}
//	return len(b), nil
//}
//
//func mockDataSrc(size int64) io.Reader {
//	fmt.Printf("dev/zero of size %d (%d MB)\n", size, size/1024/1024)
//	var z devZero
//	return io.LimitReader(z, size)
//}

func TestDecrypt(t *testing.T) {

	testfile := "/Users/gcollins/Desktop/b2test/hello1.txt"
	fo, err := os.Open(testfile)
	if err != nil {
		panic(err)
	}

	tf, err := ioutil.TempFile(os.TempDir(), "ctcloud")
	if err != nil {
		panic(err)
	}
	fmt.Println("Here is the temp: ", tf.Name())

	//keyAes, _ := hex.DecodeString(strings.Repeat("6368616e676520746869732070617373", 2))
	//keyHmac := keyAes // don't do this

	saltAR := crypt.GenSalt(nil, ScryptParams.SaltLen)
	ScryptParams.Salt = saltAR
	derivedKey, err := scrypt.Key(string("mypass"), ScryptParams)
	if err != nil {
		panic(err)
	}

	aesKey2 := derivedKey[0:16]
	//aesIv2 := derivedKey[16:32]

	err = Encrypt(fo, tf, aesKey2, aesKey2)
	if err != nil {
		log.Fatal(err)
	}
	tf.Sync()
	tf.Close()

	f := filenative.NewFile(tf.Name())

	fo, err = os.Open(f.Path())
	if err != nil {
		panic(err)
	}

	tf, err = os.Create("./hello1.txt")
	if err != nil {
		panic(err)
	}

	err = Decrypt(fo, tf, aesKey2, aesKey2)
	if err != nil {
		log.Fatal(err)
	}

}
