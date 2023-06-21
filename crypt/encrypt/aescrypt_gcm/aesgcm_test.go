package aescrypt_gcm

import (
	"strings"
	"testing"

	log "github.com/colt3k/nglog/ng"
	"github.com/colt3k/utils/encode"
	"github.com/colt3k/utils/encode/encodeenum"

	"github.com/colt3k/utils/crypt"
	"github.com/colt3k/utils/crypt/encrypt/scrypt"
)

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

func loadTestData() {
	cfg.saltScrypt = "wYr2y5H5T2/LCqWURBfVQQ=="
	cfg.saltScryptBA = []byte(cfg.saltScrypt)

	cfg.testPassword = "thisismysuperlongandcomplexpass"
	cfg.plaintext = "My Original plain text used for testing."

	cfg.cipherTxtScrypt = "l3XZaAg3JLkq/X/IFv0loI1IGAVetlmRlb1MjscWZBm9Wj9J8YEIZ2clNUALpyKrkGMD3dz51TClNaJj1E76CQ=="
	cfg.cipherTxtScryptBA = []byte(cfg.cipherTxtScrypt)
}

func TestScryptEncrypt(t *testing.T) {
	saltDecoded := encode.Decode(cfg.saltScryptBA, encodeenum.B64STD)
	log.Logln(log.DEBUG, "SaltDecoded: ", saltDecoded)
	//p, err := scrypt.Calibrate(1*time.Second, 128, scrypt.Params{})
	p := scrypt.Params{N: 65536, R: 1, P: 2, SaltLen: 16, DKLen: 32}
	log.Println("Params: ", p)

	salt := crypt.GenSalt(saltDecoded, p.SaltLen)
	log.Logf(log.DEBUG, "Salt: %s", encode.Encode(salt, encodeenum.B64STD))

	p.Salt = salt
	derivedKey, err := scrypt.Key(cfg.testPassword, p)
	if err != nil {
		panic(err)
	}
	parts := strings.Split(string(derivedKey), "$")
	lastPart := parts[len(parts)-1]
	dk := []byte(lastPart)
	log.Println("Derived Key Length : ", len(dk))

	a := New(&cfg.plaintext, nil, &dk)
	crypted := a.Encrypt()

	log.Println("CipherText: ", encode.Encode(crypted, encodeenum.B64STD))

	if encode.Encode(crypted, encodeenum.B64STD) != cfg.cipherTxtScrypt {
		t.Error("CipherText is different ", encode.Encode(crypted, encodeenum.B64STD))
	}
}

func TestScryptDecrypt(t *testing.T) {
	saltDecoded := encode.Decode(cfg.saltScryptBA, encodeenum.B64STD)
	//p, err := scrypt.Calibrate(1*time.Second, 128, scrypt.Params{})
	p := scrypt.Params{N: 65536, R: 1, P: 2, SaltLen: 16, DKLen: 32}
	log.Println("Params: ", p)

	salt := crypt.GenSalt(saltDecoded, p.SaltLen)
	log.Printf("Salt: %s\n", encode.Encode(salt, encodeenum.B64STD))

	p.Salt = salt
	derivedKey, err := scrypt.Key(cfg.testPassword, p)
	if err != nil {
		panic(err)
	}
	parts := strings.Split(string(derivedKey), "$")
	lastPart := parts[len(parts)-1]
	dk := []byte(lastPart)
	log.Println("Derived Key Length : ", len(dk))

	crypted := encode.Decode(cfg.cipherTxtScryptBA, encodeenum.B64STD)

	a := New(nil, &crypted, &dk)
	plaintext2 := a.Decrypt()

	log.Println("PlainText2: ", string(plaintext2))
	if string(plaintext2) != cfg.plaintext {
		t.Error("PlainText is different ", string(plaintext2))
	}
}
