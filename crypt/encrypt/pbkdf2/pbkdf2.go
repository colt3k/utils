package pbkdf2

import (
	"crypto/sha1"

	"golang.org/x/crypto/pbkdf2"

	"github.com/colt3k/utils/crypt"
)

/*
SECRET_KEY_ALGORITHM 	: PBKDF2WithHmacSHA1

PBKDF2WithHmacSHA256
PBKDF2: key derivation function PBKDF2 as defined in RFC2898 / PKCS #5 v2.0
i.e. PBKDF2 SHA1 vs SHA256 Hash algorithm strength is important, but it is not so important in key derivation functions.
It is unlikely that even if SHA-1 is broken that it would influence the security of PBKDF2. You are better off using
SHA-1, and increase the iteration count up to a level that is tweaked for your specific configuration.
If you want to protect against hardware acceleration use SCrypt Instead of PBKDF2

Password is passed by user
Salt is a unique salt for the system this will be running on and doesn't change
*/

const (
	iterationsDFLT = 65536
	keyLengthDFLT  = 32
)

type PBKDF2 struct {
	pass       []byte
	salt       []byte
	keyLength  int
	iterations int
}

func New(pass, salt []byte, keyLength, iterations int) *PBKDF2 {
	t := new(PBKDF2)
	t.pass = pass
	t.salt = salt
	t.keyLength = keyLengthDFLT
	if keyLength > 0 {
		t.keyLength = keyLength
	}
	t.iterations = iterationsDFLT
	if iterations > 0 {
		t.iterations = iterations
	}

	return t
}

func (p *PBKDF2) Generate() []byte {
	p.salt = crypt.GenSalt(p.salt, p.keyLength/2)
	return pbkdf2.Key(p.pass, p.salt, p.iterations, p.keyLength, sha1.New)
}
