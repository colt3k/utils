package scrypt

import (
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/scrypt"

	log "github.com/colt3k/nglog/ng"
	"github.com/colt3k/utils/crypt"
)

var (
	// double of the SALT
	aesKeyLength = 32
	// half of the key length
	//saltLength = aesKeyLength / 2

	errSalt = errors.New("salt is nil")
)

// Constants
const (
	maxInt     = 1<<31 - 1
	minDKLen   = 16 // the minimum derived key length in bytes.
	minSaltLen = 8  // the minimum allowed salt length in bytes.
)

// ErrInvalidParams is returned when the cost parameters (N, r, p), salt length
// or derived key length are invalid.
var ErrInvalidParams = errors.New("scrypt: the parameters provided are invalid")

// ErrInvalidHash is returned when failing to parse a provided scrypt
// hash and/or parameters.
var ErrInvalidHash = errors.New("scrypt: the provided hash is not in the correct format")

// ErrMismatchedHashAndPassword is returned when a password (hashed) and
// given hash do not match.
var ErrMismatchedHashAndPassword = errors.New("scrypt: the hashed password does not match the hash of the given password")

// DefaultParams provides sensible default inputs into the scrypt function
// for interactive use (i.e. web applications).
// These defaults will consume approxmiately 16MB of memory (128 * r * N).
// The default key length is 256 bits.
var DefaultParams = Params{N: 16384, R: 8, P: 1, SaltLen: 16, DKLen: aesKeyLength}

/*
Key derived key for e.g. AES-256
Key derives a key from the password, salt, and cost parameters returning
	a byte slice of length keyLen that can be used as cryptographic key.

N is a CPU/memory cost parameter, which must be a power of two greater than 1.

r and p must satisfy r * p < 2³⁰. If the parameters do not satisfy the limits,
	the function returns a nil byte slice and an error.

Password is passed by user
Salt is a unique salt for the system this will be running on and doesn't change
*/
func Key(password string, params Params) ([]byte, error) {

	if params.Salt == nil {
		return nil, errSalt
	}

	if err := params.Check(); err != nil {
		return nil, err
	}
	// They should be increased as memory latency and CPU parallelism increases. As OF 2009

	//This would create an AES256 Key of 32bytes:
	// 		dk, err := scrypt.Key([]byte("some password"), salt, 16384, 8, 1, 32)
	// Interactive Logins should use at least as of 2017 are N=32768, r=8 and p=1.
	dk, err := scrypt.Key([]byte(password), params.Salt, params.N, params.R, params.P, params.DKLen)
	if err != nil {
		return nil, err
	}

	// Prepend the params and the salt to the derived key, each separated
	// by a "$" character. The salt and the derived key are hex encoded.
	//fmt.Println("SALT IN SCRYPT: ", Base64encode(params.Salt))
	//log.Logln(log.DEBUG, "Building Derived Key, pass back in format N$R$P$SALT$DK")
	return []byte(fmt.Sprintf("%d$%d$%d$%x$%x", params.N, params.R, params.P, params.Salt, dk)), nil
}

// CompareHashAndPassword compares a derived key with the possible cleartext
// equivalent. The parameters used in the provided derived key are used.
// The comparison performed by this function is constant-time. It returns nil
// on success, and an error if the derived keys do not match.
func CompareHashAndPassword(hash []byte, password []byte) error {
	// Decode existing hash, retrieve params and salt.
	params, salt, dk, err := decodeHash(hash)
	if err != nil {
		return err
	}

	// scrypt the cleartext password with the same parameters and salt
	other, err := scrypt.Key(password, salt, params.N, params.R, params.P, params.DKLen)
	if err != nil {
		return err
	}

	// Constant time comparison
	if subtle.ConstantTimeCompare(dk, other) == 1 {
		return nil
	}

	return ErrMismatchedHashAndPassword
}

// decodeHash extracts the parameters, salt and derived key from the
// provided hash. It returns an error if the hash format is invalid and/or
// the parameters are invalid.
func decodeHash(hash []byte) (Params, []byte, []byte, error) {
	vals := strings.Split(string(hash), "$")

	log.Println("DecodeHash len: ", len(vals))
	// P, N, R, salt, scrypt derived key
	if len(vals) != 5 {
		return Params{}, nil, nil, ErrInvalidHash
	}

	var params Params
	var err error

	params.N, err = strconv.Atoi(vals[0])
	if err != nil {
		return params, nil, nil, ErrInvalidHash
	}

	params.R, err = strconv.Atoi(vals[1])
	if err != nil {
		return params, nil, nil, ErrInvalidHash
	}

	params.P, err = strconv.Atoi(vals[2])
	if err != nil {
		return params, nil, nil, ErrInvalidHash
	}

	salt, err := hex.DecodeString(vals[3])
	if err != nil {
		return params, nil, nil, ErrInvalidHash
	}
	params.SaltLen = len(salt)

	dk, err := hex.DecodeString(vals[4])
	if err != nil {
		return params, nil, nil, ErrInvalidHash
	}
	params.DKLen = len(dk)

	if err := params.Check(); err != nil {
		return params, nil, nil, err
	}

	log.Println("To normal return")
	return params, salt, dk, nil
}

// Cost returns the scrypt parameters used to generate the derived key. This
// allows a package user to increase the cost (in time & resources) used as
// computational performance increases over time.
func Cost(hash []byte) (Params, error) {
	params, _, _, err := decodeHash(hash)

	return params, err
}

// Calibrate returns the hardest parameters (not weaker than the given params),
// allowed by the given limits.
// The returned params will not use more memory than the given (MiB);
// will not take more time than the given timeout, but more than timeout/2.
//
//   The default timeout (when the timeout arg is zero) is 200ms.
//   The default memMiBytes (when memMiBytes is zero) is 16MiB.
//   The default parameters (when params == Params{}) is DefaultParams.
func Calibrate(timeout time.Duration, memMiBytes int, params Params) (Params, error) {
	p := params
	if p.N == 0 || p.R == 0 || p.P == 0 || p.SaltLen == 0 || p.DKLen == 0 {
		log.Logln(log.DEBUG, "using default scrypt params N(iterations): 16384, R(blocksize): 8, P(parallel): 1, SaltLen: 16, DerivedKeyLen: 32")
		p = DefaultParams
	} else if err := p.Check(); err != nil {
		return p, err
	}
	if timeout == 0 {
		timeout = 200 * time.Millisecond
	}
	if memMiBytes == 0 {
		memMiBytes = 16
	}
	log.Logln(log.DEBUG, "generating salt of length", p.SaltLen)
	salt := crypt.GenerateRandomBytes(p.SaltLen)

	password := []byte("weakpassword")

	// First, we calculate the minimal required time.
	log.Logln(log.DEBUG, "calculating the minimal required time")
	start := time.Now()
	if _, err := scrypt.Key(password, salt, p.N, p.R, p.P, p.DKLen); err != nil {
		return p, err
	}
	dur := time.Since(start)

	for dur < timeout && p.N < maxInt>>1 {
		p.N <<= 1
	}

	// Memory usage is at least 128 * r * N, see
	// http://blog.ircmaxell.com/2014/03/why-i-dont-recommend-scrypt.html
	// or
	// https://drupal.org/comment/4675994#comment-4675994

	var again bool
	memBytes := memMiBytes << 20
	log.Logln(log.DEBUG, "compute amount of memory used is less than ", memMiBytes, "shifted left", memBytes)
	// If we'd use more memory then the allowed, we can tune the memory usage
	for 128*int64(p.R)*int64(p.N) > int64(memBytes) {
		if p.R > 1 {
			// by lowering r
			p.R--
		} else if p.N > 16 {
			again = true
			p.N >>= 1
		} else {
			break
		}
	}
	if !again {
		return p, p.Check()
	}

	// We have to compensate the lowering of N, by increasing p.
	log.Logln(log.DEBUG, "compensating the lowering of N by increasing p, N:", p.N, "P:", p.P)
	for i := 0; i < 10 && p.P > 0; i++ {
		start := time.Now()
		if _, err := scrypt.Key(password, salt, p.N, p.R, p.P, p.DKLen); err != nil {
			return p, err
		}
		dur := time.Since(start)
		if dur < timeout/2 {
			p.P = int(float64(p.P)*float64(timeout/dur) + 1)
		} else if dur > timeout && p.P > 1 {
			p.P--
		} else {
			break
		}
	}

	return p, p.Check()
}

// Check checks that the parameters are valid for input into the
// scrypt key derivation function.
func (p *Params) Check() error {
	// Validate N
	//log.Logln(log.DEBUG, "VALIDATE N (iterations): less than maxInt OR less than equal to 1 OR modulo 2 not equal zero")
	if p.N > maxInt || p.N <= 1 || p.N%2 != 0 {
		return ErrInvalidParams
	}

	//log.Logln(log.DEBUG, "VALIDATE R (block size): less than 1 or greater than maxInt")
	// Validate r
	if p.R < 1 || p.R > maxInt {
		return ErrInvalidParams
	}

	//log.Logln(log.DEBUG, "VALIDATE P (parallelism): less than one OR greater than maxInt")
	// Validate p
	if p.P < 1 || p.P > maxInt {
		return ErrInvalidParams
	}

	// Validate that r & p don't exceed 2^30 and that N, r, p values don't
	// exceed the limits defined by the scrypt algorithm.
	if uint64(p.R)*uint64(p.P) >= 1<<30 || p.R > maxInt/128/p.P || p.R > maxInt/256 || p.N > maxInt/128/p.R {
		return ErrInvalidParams
	}

	//log.Logln(log.DEBUG, "VALIDATE (salt len): less than minimum 8 OR greater than maxInt")
	// Validate the salt length
	if p.SaltLen < minSaltLen || p.SaltLen > maxInt {
		return ErrInvalidParams
	}

	//log.Logln(log.DEBUG, "VALIDATE (derived key length): less than 16 OR greater than maxInt")
	// Validate the derived key length
	if p.DKLen < minDKLen || p.DKLen > maxInt {
		return ErrInvalidParams
	}

	return nil
}
