package aescrypt

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
	"time"

	log "github.com/colt3k/nglog/ng"
	"github.com/colt3k/utils/encode"
	"github.com/colt3k/utils/encode/encodeenum"
	"golang.org/x/crypto/poly1305"

	"github.com/colt3k/utils/crypt"
	"github.com/colt3k/utils/crypt/encrypt/pbkdf2"
	"github.com/colt3k/utils/crypt/encrypt/scrypt"
)

// ALLOW UPDATE OF OUTPUT WITH FLAG pass -update
var update = flag.Bool("update", false, "update .golden files")

type appConfig struct {
	salt              string
	saltBA            []byte
	saltScrypt        string
	saltScryptBA      []byte
	testPassword      string
	plaintext         string
	cipherTxt         string
	cipherTxtBA       []byte
	cipherTxtScrypt   string
	cipherTxtScryptBA []byte
}

var cfg *appConfig

func init() {
	log.SetLevel(log.DEBUG)
	cfg = &appConfig{}
	loadTestData()
}

//**************************** END OF OUTPUT WITH FLAG

func loadTestData() {

	//Pull / Convert Test Data
	cfg.salt = "qhebZTd7PVqGBCH0rTyl0w=="
	cfg.saltBA = []byte(cfg.salt)

	cfg.saltScrypt = "wYr2y5H5T2/LCqWURBfVQQ=="
	cfg.saltScryptBA = []byte(cfg.saltScrypt)

	cfg.testPassword = "thisismysuperlongandcomplexpass"

	cfg.plaintext = "My Original plain text used for testing."

	cfg.cipherTxt = "e42ax9c+EW8uIZB+uNQdIXysO5w6YWkl6TcG67L6XM7dyvtzwo3FhKiJ30C/Qxub"
	cfg.cipherTxtBA = []byte(cfg.cipherTxt)

	cfg.cipherTxtScrypt = "E7CqSdOc1pWZwPEw60eM4OiNMc7vxxr5l6wGwQzNdqiIXIWU84VQyBj0ZvhDaG5r"
	cfg.cipherTxtScryptBA = []byte(cfg.cipherTxtScrypt)
}

func buildOutput(t *testing.T, name string, actual []byte) {

	golden := filepath.Join("testdata", name+".golden")
	//If our update flag is on then write it
	if *update {
		ioutil.WriteFile(golden, actual, 0644)
	}
	//Read our golden data
	expected, _ := ioutil.ReadFile(golden)

	// FAIL!
	if !bytes.Equal(actual, expected) {
		t.Error("Output doesn't match.")
	}
}

func TestEncrypt(t *testing.T) {

	var buff bytes.Buffer

	log.Logln(log.INFO, "Update Output file?", *update)

	bytAr := []byte(cfg.salt)
	saltDecoded := encode.Decode(bytAr, encodeenum.B64STD)
	log.Logln(log.DEBUG, "SaltDecoded: ", saltDecoded)

	saltAR := crypt.GenSalt(saltDecoded, ivSize)
	log.Logln(log.DEBUG, "Salt: ", saltAR)

	encSalt := encode.Encode(saltAR, encodeenum.B64STD)
	buff.WriteString(fmt.Sprintf("Salt: %s\n", encSalt))

	derivedKey := pbkdf2.New([]byte(cfg.testPassword), saltAR, aesKeyLength, iterationCount).Generate()
	buff.WriteString(fmt.Sprintf("Derived Key Length: %d\n", len(derivedKey)))

	a := New(-1, -1, -1, -1, &cfg.plaintext, nil, &derivedKey)
	crypted := a.Encrypt()
	buff.WriteString(fmt.Sprintf("CipherText: %s", encode.Encode(crypted, encodeenum.B64STD)))

	buildOutput(t, "encrypt", buff.Bytes())

}
func TestDecrypt(t *testing.T) {

	var buff bytes.Buffer

	saltDecoded := encode.Decode(cfg.saltBA, encodeenum.B64STD)

	salt := crypt.GenSalt(saltDecoded, saltLength)
	buff.WriteString(fmt.Sprintf("Salt: %s\n", encode.Encode(salt, encodeenum.B64STD)))

	derivedKey := pbkdf2.New([]byte(cfg.testPassword), salt, aesKeyLength, iterationCount).Generate()
	buff.WriteString(fmt.Sprintf("Derived Key Length: %d\n", len(derivedKey)))

	crypted := encode.Decode(cfg.cipherTxtBA, encodeenum.B64STD)

	a := New(-1, -1, -1, -1, nil, &crypted, &derivedKey)
	plaintext2 := a.Decrypt()

	buff.WriteString(fmt.Sprintf("PlainText2: %s", strings.TrimSpace(string(plaintext2))))

	buildOutput(t, "decrypt", buff.Bytes())
}

func TestScryptEncrypt(t *testing.T) {

	var buf bytes.Buffer

	saltDecoded := encode.Decode(cfg.saltScryptBA, encodeenum.B64STD)

	//Calibrate to current system
	p, err := scrypt.Calibrate(1*time.Second, 128, scrypt.Params{})

	if d, err := json.MarshalIndent(p, "", "    "); err == nil {
		buf.WriteString(fmt.Sprintln("Calibration - ScryptParams: ", string(d)))
	}

	//Generate salt if required
	salt := crypt.GenSalt(saltDecoded, p.SaltLen)
	buf.WriteString(fmt.Sprintf("Salt: %s\n", encode.Encode(salt, encodeenum.B64STD)))

	//Store salt in SCryptParams
	p.Salt = salt
	//Derive a key from SCryptParams
	derivedKey, err := scrypt.Key(cfg.testPassword, p)
	if err != nil {
		panic(err)
	}

	buf.WriteString(fmt.Sprintf("Derived Key Length (inludes SCryptParam data): %d\n", len(derivedKey)))

	//Encrypt plaintext using derived Key
	a := New(-1, -1, -1, -1, &cfg.plaintext, nil, &derivedKey)
	crypted := a.Encrypt()

	// **************************** MAC CREATION AND VALIDATION -- START **************
	var out [16]byte
	var k [32]byte
	//Copy created key into k slice of 32
	copy(k[:16], derivedKey[:])
	//Pass in empty slice of 16 for mac key storage, encrypted text and derived Key
	poly1305.Sum(&out, crypted, &k)

	//Convert our slice into a []byte
	macKey := make([]byte, len(k))
	copy(macKey, out[:])
	//Encode our []byte macKey for storage
	macKeyStr := encode.Encode(macKey, encodeenum.B64STD)
	buf.WriteString(fmt.Sprintln("Mac: ", macKeyStr))

	//Convert Mac Key back for use as k slice
	//base 64 decode our macKey String
	macKeyBA := []byte(macKeyStr)
	macKeyDecoded := encode.Decode(macKeyBA, encodeenum.B64STD)

	//Copy macKey decoded string into out slice
	copy(out[:], []byte(macKeyDecoded))

	//Verify our security is still in place
	buf.WriteString(fmt.Sprintln("Verifying Security of data..."))
	if !poly1305.Verify(&out, crypted, &k) {
		buf.WriteString(fmt.Sprintln("NOT Valid"))
	}
	// **************************** MAC CREATION AND VALIDATION -- END ******************
	buf.WriteString(fmt.Sprintf("CipherText: %s", encode.Encode(crypted, encodeenum.B64STD)))

	buildOutput(t, "scryptencrypt", buf.Bytes())
}
func TestScryptDecrypt(t *testing.T) {

	var buf bytes.Buffer

	saltDecoded := encode.Decode(cfg.saltScryptBA, encodeenum.B64STD)

	//Calibrate to current system
	p, err := scrypt.Calibrate(1*time.Second, 128, scrypt.Params{})
	if d, err := json.MarshalIndent(p, "", "    "); err == nil {
		buf.WriteString(fmt.Sprintln("Calibration - ScryptParams: ", string(d)))
	}

	//Generate salt if required
	salt := crypt.GenSalt(saltDecoded, p.SaltLen)
	buf.WriteString(fmt.Sprintf("Salt: %s\n", encode.Encode(salt, encodeenum.B64STD)))

	//Store salt in SCryptParams
	p.Salt = salt
	//Derive a key from SCryptParams
	derivedKey, err := scrypt.Key(cfg.testPassword, p)
	if err != nil {
		panic(err)
	}
	buf.WriteString(fmt.Sprintf("Derived Key Length (inludes SCryptParam data): %d\n", len(derivedKey)))

	//Decode Cipher Text
	crypted := encode.Decode(cfg.cipherTxtScryptBA, encodeenum.B64STD)

	//Decrypt ciphertext using derived Key
	a := New(-1, -1, -1, -1, nil, &crypted, &derivedKey)
	plaintext2 := a.Decrypt()

	buf.WriteString(fmt.Sprintf("PlainText2: %s", strings.TrimSpace(string(plaintext2))))

	buildOutput(t, "scryptdecrypt", buf.Bytes())
}

func BenchmarkAESPBKDF2Crypt(b *testing.B) {

	for i := 0; i < b.N; i++ {
		saltDecoded := encode.Decode(cfg.saltBA, encodeenum.B64STD)

		salt := crypt.GenSalt(saltDecoded, saltLength)
		log.Println("Salt: ", encode.Encode(salt, encodeenum.B64STD))

		derivedKey := pbkdf2.New([]byte(cfg.testPassword), salt, aesKeyLength, iterationCount).Generate()

		log.Println("Derived Key Length : ", len(derivedKey))

		a := New(-1, -1, -1, -1, &cfg.plaintext, nil, &derivedKey)
		crypted := a.Encrypt()

		log.Println("CipherText: ", encode.Encode(crypted, encodeenum.B64STD))
	}
}

func BenchmarkAESDecrypt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		saltDecoded := encode.Decode(cfg.saltBA, encodeenum.B64STD)
		salt := crypt.GenSalt(saltDecoded, saltLength)
		log.Println("Salt: ", encode.Encode(salt, encodeenum.B64STD))
		derivedKey := pbkdf2.New([]byte(cfg.testPassword), salt, aesKeyLength, iterationCount).Generate()
		//derivedKey := PBKDF2Key(&cfg.testPassword, salt)

		log.Println("Derived Key Length : ", len(derivedKey))

		crypted := encode.Decode(cfg.cipherTxtBA, encodeenum.B64STD)

		a := New(-1, -1, -1, -1, nil, &crypted, &derivedKey)
		plaintext2 := a.Decrypt()

		log.Println("PlainText2: ", strings.TrimSpace(string(plaintext2)))
	}
}

func BenchmarkAESScryptCrypt(b *testing.B) {

	for i := 0; i < b.N; i++ {
		saltDecoded := encode.Decode(cfg.saltScryptBA, encodeenum.B64STD)
		p, err := scrypt.Calibrate(1*time.Second, 128, scrypt.Params{})
		log.Println("Params: ", p)
		salt := crypt.GenSalt(saltDecoded, p.SaltLen)
		log.Printf("Salt: %s\n", encode.Encode(salt, encodeenum.B64STD))
		p.Salt = salt
		derivedKey, err := scrypt.Key(cfg.testPassword, p)
		if err != nil {
			panic(err)
		}

		log.Println("Derived Key Length : ", len(derivedKey))
		a := New(-1, -1, -1, -1, &cfg.plaintext, nil, &derivedKey)
		crypted := a.Encrypt()

		log.Println("CipherText: ", encode.Encode(crypted, encodeenum.B64STD))
	}
}

func BenchmarkAESScryptDecrypt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		saltDecoded := encode.Decode(cfg.saltScryptBA, encodeenum.B64STD)
		p, err := scrypt.Calibrate(1*time.Second, 128, scrypt.Params{})
		log.Println("Params: ", p)
		salt := crypt.GenSalt(saltDecoded, p.SaltLen)
		log.Printf("Salt: %s\n", encode.Encode(salt, encodeenum.B64STD))
		p.Salt = salt
		derivedKey, err := scrypt.Key(cfg.testPassword, p)
		if err != nil {
			panic(err)
		}

		log.Println("Derived Key Length : ", len(derivedKey))

		crypted := encode.Decode(cfg.cipherTxtScryptBA, encodeenum.B64STD)

		a := New(-1, -1, -1, -1, nil, &crypted, &derivedKey)
		plaintext2 := a.Decrypt()

		log.Println("PlainText2: ", strings.TrimSpace(string(plaintext2)))
	}
}

func TestAESCrypt(t *testing.T) {

	saltDecoded := encode.Decode(cfg.saltScryptBA, encodeenum.B64STD)
	p, err := scrypt.Calibrate(1*time.Second, 128, scrypt.Params{})
	log.Println("Params: ", p)
	salt := crypt.GenSalt(saltDecoded, p.SaltLen)
	log.Printf("Salt: %s\n", encode.Encode(salt, encodeenum.B64STD))
	p.Salt = salt
	derivedKey, err := scrypt.Key(cfg.testPassword, p)
	if err != nil {
		panic(err)
	}

	log.Println("Derived Key Length : ", len(derivedKey))

	a := New(-1, -1, -1, -1, &cfg.plaintext, nil, &derivedKey)
	crypted := a.Encrypt()

	log.Println("CipherText: ", encode.Encode(crypted, encodeenum.B64STD))

	if encode.Encode(crypted, encodeenum.B64STD) != cfg.cipherTxtScrypt {
		t.Error("CipherText is different ", encode.Encode(crypted, encodeenum.B64STD))
	}
}
func TestAESDecrypt(t *testing.T) {
	saltDecoded := encode.Decode(cfg.saltScryptBA, encodeenum.B64STD)
	p, err := scrypt.Calibrate(1*time.Second, 128, scrypt.Params{})
	log.Println("Params: ", p)
	salt := crypt.GenSalt(saltDecoded, p.SaltLen)
	log.Printf("Salt: %s\n", encode.Encode(salt, encodeenum.B64STD))
	p.Salt = salt
	derivedKey, err := scrypt.Key(cfg.testPassword, p)
	if err != nil {
		panic(err)
	}

	log.Println("Derived Key Length : ", len(derivedKey))

	crypted := encode.Decode(cfg.cipherTxtScryptBA, encodeenum.B64STD)

	a := New(-1, -1, -1, -1, nil, &crypted, &derivedKey)
	plaintext2 := a.Decrypt()

	log.Println("PlainText2: ", string(plaintext2))
	if string(plaintext2) != cfg.plaintext {
		t.Error("PlainText is different ", string(plaintext2))
	}
}
