package sign

/*
https://github.com/restic/restic/tree/master/internal/crypto
https://restic.readthedocs.io/en/latest/100_references.html

*/
import (
	mrand "math/rand"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	log "github.com/colt3k/nglog/ng"

	"github.com/colt3k/utils/crypt"
)

var exText = "My Example Text"

// Random returns size bytes of pseudo-random data derived from the seed.
func Random(seed, count int) []byte {
	p := make([]byte, count)

	rnd := mrand.New(mrand.NewSource(int64(seed)))

	for i := 0; i < len(p); i += 8 {
		val := rnd.Int63()
		var data = []byte{
			byte((val >> 0) & 0xff),
			byte((val >> 8) & 0xff),
			byte((val >> 16) & 0xff),
			byte((val >> 24) & 0xff),
			byte((val >> 32) & 0xff),
			byte((val >> 40) & 0xff),
			byte((val >> 48) & 0xff),
			byte((val >> 56) & 0xff),
		}

		for j := range data {
			cur := i + j
			if cur >= len(p) {
				break
			}
			p[cur] = data[j]
		}
	}

	return p
}

// Assert fails the test if the condition is false.
func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		log.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// OK fails the test if an err is not nil.
func OK(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		log.Printf("\033[31m%s:%d: unexpected error: %+v\033[39m\n\n", filepath.Base(file), line, err)
		tb.FailNow()
	}
}

// Equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		log.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func TestMine(t *testing.T) {
	//Generate a Random Key
	k := NewRandomKey()

	//Create a buffer of the data size + Extension(ivSize+macSize) == (aes.BlockSize)+(poly1305.TagSize)
	buf := make([]byte, 0, len(exText)+Extension)

	//Create Random Nonce AKA IV (initialization vector), saved with Key for storage
	nonce := crypt.GenerateRandomBytes(ivSize)

	//Seal our data pass in our dst, iv, plaintext, nil for additional Data
	//MAC is stored in buf
	ciphertext := k.Seal(buf[:0], nonce, []byte(exText), nil)

	trusted, err := k.Safe(nonce, ciphertext)
	if err != nil {
		log.Printf("Error: ", err)
		t.FailNow()
	}
	if !trusted {
		println("Cipher text has been tampered with and isn't to be trusted.")
	}
	println("Cipher text has NOT been tampered with and CAN be trusted.")

	//Create byte array the size of our ciphertext
	plaintext := make([]byte, 0, len(ciphertext))

	//Decrypt and Authenticate via Open pass in dst, IV, ciphertext, nil for additional data
	plaintext, err = k.Open(plaintext[:0], nonce, ciphertext, nil)
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		log.Printf("\033[31m%s:%d: unexpected error: %+v\033[39m\n\n", filepath.Base(file), line, err)
		t.FailNow()
	}
	log.Printf("PlainText: |%s|\n|Original: %s|\n", plaintext, exText)

	log.Println("Key: ", k)
}

func TestSomething(t *testing.T) {

	//Generate a Random Key
	k := NewRandomKey()

	tests := []int{5, 23, 2<<18 + 23, 1 << 20}

	for _, size := range tests {
		//Generate a random set of data
		data := Random(42, size)

		//Create a buffer of the data size + Extension(ivSize+macSize) == (aes.BlockSize)+(poly1305.TagSize)
		buf := make([]byte, 0, size+Extension)

		//Create Random Nonce AKA IV (initialization vector)
		nonce := NewRandomNonce()
		//Seal our data pass in our dst, iv, plaintext, nil for additional Data
		ciphertext := k.Seal(buf[:0], nonce, data, nil)

		//Length of ciphertext is equal to length of data plus poly1305.TagSize(Overhead)
		Assert(t, len(ciphertext) == len(data)+k.Overhead(),
			"ciphertext length does not match: want %d, got %d",
			len(data)+Extension, len(ciphertext))

		//Create byte array the size of our ciphertext
		plaintext := make([]byte, 0, len(ciphertext))

		//Decrypt and Authenticate via Open pass in dst, IV, ciphertext, nil for additional data
		plaintext, err := k.Open(plaintext[:0], nonce, ciphertext, nil)
		//Ensure err is nil
		OK(t, err)

		//CHECK If the length of plaintext is equal to length of original data
		Assert(t, len(plaintext) == len(data),
			"plaintext length does not match: want %d, got %d",
			len(data), len(plaintext))

		//Does the plaintext equal the original data?
		Equals(t, plaintext, data)

	}
	log.Println("Key: ", k)
}
