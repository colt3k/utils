package sign

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/crypto/poly1305"

	log "github.com/colt3k/nglog/ng"
)

const (
	AESKeyLength = 32 // double of the SALT
	//Poly1305 data signing (MAC Message Authentication Code)
	macKeySizeK = 16 // for AES-128
	macKeySizeR = 16 // for Poly1305
	//macKeySize  = macKeySizeK + macKeySizeR // for Poly1305-AES128
	ivSize = aes.BlockSize

	macSize = poly1305.TagSize
	// Extension is the number of bytes a plaintext is enlarged by encrypting it.
	Extension = ivSize + macSize
)

var (

	//ErrUnauthenticated is returned when ciphertext verification has failed.
	ErrUnauthenticated = errors.New("ciphertext verification failed")
	//ErrSalt salt error message
	//ErrSalt = errors.New("salt is nil")
)

// statically ensure that *Key implements crypto/cipher.AEAD
//var _ cipher.AEAD = &Key{}

// mask for key, (cf. http://cr.yp.to/mac/poly1305-20050329.pdf)
var poly1305KeyMask = [16]byte{
	0xff,
	0xff,
	0xff,
	0x0f, // 3: top four bits zero
	0xfc, // 4: bottom two bits zero
	0xff,
	0xff,
	0x0f, // 7: top four bits zero
	0xfc, // 8: bottom two bits zero
	0xff,
	0xff,
	0x0f, // 11: top four bits zero
	0xfc, // 12: bottom two bits zero
	0xff,
	0xff,
	0x0f, // 15: top four bits zero
}

func poly1305MAC(msg []byte, nonce []byte, key *MACKey) []byte {
	k := poly1305PrepareKey(nonce, key)

	var out [16]byte
	poly1305.Sum(&out, msg, &k)

	return out[:]
}

// mask poly1305 key
func maskKey(k *MACKey) {
	if k == nil || k.masked {
		return
	}

	for i := 0; i < poly1305.TagSize; i++ {
		k.R[i] = k.R[i] & poly1305KeyMask[i]
	}

	k.masked = true
}

// prepare key for low-level poly1305.Sum(): r||n
func poly1305PrepareKey(nonce []byte, key *MACKey) [32]byte {
	var k [32]byte

	maskKey(key)

	//Create an AES Cipher
	aesCipher, err := aes.NewCipher(key.K[:])
	if err != nil {
		panic(err)
	}
	//Encrypt nonce (all of it since nonce is a slice)
	aesCipher.Encrypt(k[16:], nonce[:])

	//Copy key.R (MACKey) to k 0 to 16 of slice
	copy(k[:16], key.R[:])

	return k
}

func poly1305Verify(msg []byte, nonce []byte, key *MACKey, mac []byte) bool {
	k := poly1305PrepareKey(nonce, key)

	var m [16]byte
	copy(m[:], mac)

	return poly1305.Verify(&m, msg, &k)
}

// NewRandomKey returns new encryption and message authentication keys.
func NewRandomKey() *Key {
	//Create empty struct
	k := &Key{}

	n, err := rand.Read(k.EncryptionKey[:])
	if n != AESKeyLength || err != nil {
		panic("unable to read enough random bytes for encryption key")
	}

	n, err = rand.Read(k.MACKey.K[:])
	if n != macKeySizeK || err != nil {
		panic("unable to read enough random bytes for MAC encryption key")
	}

	n, err = rand.Read(k.MACKey.R[:])
	if n != macKeySizeR || err != nil {
		panic("unable to read enough random bytes for MAC key")
	}

	maskKey(&k.MACKey)
	return k
}

// NewRandomNonce returns a new random nonce. It panics on error so that the
// program is safely terminated.
func NewRandomNonce() []byte {
	iv := make([]byte, ivSize)
	n, err := rand.Read(iv)
	if n != ivSize || err != nil {
		panic("unable to read enough random bytes for iv")
	}
	return iv
}

// Seal encrypts and authenticates plaintext, authenticates the
// additional data and appends the result to dst, returning the updated
// slice. The nonce must be NonceSize() bytes long and unique for all
// time, for a given key.
//
// The plaintext and dst may alias exactly or not at all. To reuse
// plaintext's storage for the encrypted output, use plaintext[:0] as dst.
func (k *Key) Seal(dst, nonce, plaintext, additionalData []byte) []byte {
	/* ***************** VALIDATION *******************/
	if !k.Valid() {
		panic("key is invalid")
	}

	if len(additionalData) > 0 {
		panic("additional data is not supported")
	}

	if len(nonce) != ivSize {
		panic("incorrect nonce length")
	}

	if !validNonce(nonce) {
		panic("nonce is invalid")
	}

	//Create slice 'ret' (HEAD) and 'out' (TAIL)
	log.Println("DST b4 slice4Append EMPTY: ", dst)
	log.Println("Size of slice4Append 31: ", len(plaintext)+k.Overhead())
	log.Println("Size of slice4Append PLAIN: ", len(plaintext))
	ret, out := sliceForAppend(dst, len(plaintext)+k.Overhead())
	log.Println("POST slice4Append HEAD (zerod): ", ret)
	log.Println("POST slice4Append TAIL (zerod): ", out)

	//START ENCRYPTION *********************************

	//CTR converts a block cipher into a stream cipher by
	// repeatedly encrypting an incrementing counter and
	// xoring the resulting stream of data with the input.
	c, err := aes.NewCipher(k.EncryptionKey[:])
	if err != nil {
		panic(fmt.Sprintf("unable to create cipher: %v", err))
	}
	//Pass Cipher to CTR (creating a stream encryptor)
	e := cipher.NewCTR(c, nonce)
	//Pass all data and run through CTR XOR as encrypted to OUT
	e.XORKeyStream(out, plaintext)

	//END ENCRYPTION ************************************

	//********* MAC GENERATION *************************
	// pass in TAIL, NONCE, MACKEY
	log.Print("\n\n")
	log.Println("B4 Poly1305MAC TAIL : ", out[:len(plaintext)])
	log.Println("B4 Poly1305MAC NONCE : ", nonce)
	mval, iss := k.MACKey.MarshalJSON()
	if iss != nil {
		panic(iss)
	}

	log.Println("B4 Poly1305MAC MACKEY : ", string(mval))
	mac := poly1305MAC(out[:len(plaintext)], nonce, &k.MACKey)
	copy(out[len(plaintext):], mac)

	return ret
}

/*
sliceForAppend takes a slice and a requested number of bytes. It returns a
slice with the contents of the given slice followed by that many bytes and a
second slice that aliases into it and contains only the extra bytes. If the
original slice has sufficient capacity then no allocation is performed.

taken from the stdlib, crypto/aes/aes_gcm.go
*/
func sliceForAppend(in []byte, n int) (head, tail []byte) {
	if total := len(in) + n; cap(in) >= total {
		head = in[:total]
	} else {
		head = make([]byte, total)
		copy(head, in)
	}
	tail = head[len(in):]
	return
}

// validNonce checks that nonce is not all zero.
func validNonce(nonce []byte) bool {
	var sum byte
	for _, b := range nonce {
		sum |= b
	}
	return sum > 0
}

/*
 Overhead returns the maximum difference between the lengths of a
 plaintext and its ciphertext.
*/
func (k *Key) Overhead() int {
	return macSize
}

// Valid tests if the key is valid.
func (k *Key) Valid() bool {
	return k.EncryptionKey.Valid() && k.MACKey.Valid()
}

// Valid tests whether the key k is valid (i.e. not zero).
func (k *EncryptionKey) Valid() bool {
	for i := 0; i < len(k); i++ {
		if k[i] != 0 {
			return true
		}
	}

	return false
}

// Valid tests whether the key k is valid (i.e. not zero).
func (m *MACKey) Valid() bool {
	nonzeroK := false
	for i := 0; i < len(m.K); i++ {
		if m.K[i] != 0 {
			nonzeroK = true
		}
	}

	if !nonzeroK {
		return false
	}

	for i := 0; i < len(m.R); i++ {
		if m.R[i] != 0 {
			return true
		}
	}

	return false
}

/*
Open decrypts and authenticates ciphertext, authenticates the
additional data and, if successful, appends the resulting plaintext
to dst, returning the updated slice. The nonce must be NonceSize()
bytes long and both it and the additional data must match the
value passed to Seal.

The ciphertext and dst may alias exactly or not at all. To reuse
ciphertext's storage for the decrypted output, use ciphertext[:0] as dst.

Even if the function fails, the contents of dst, up to its capacity,
may be overwritten.
*/
func (k *Key) Open(dst, nonce, ciphertext, additionalData []byte) ([]byte, error) {
	if !k.Valid() {
		return nil, errors.New("invalid key")
	}

	// check parameters
	if len(nonce) != ivSize {
		panic("incorrect nonce length")
	}

	if !validNonce(nonce) {
		return nil, errors.New("nonce is invalid")
	}

	// check for plausible length
	if len(ciphertext) < k.Overhead() {
		return nil, errors.New("trying to decrypt invalid data: ciphertext too small")
	}

	//Size is length of ciphertext minus macSize
	l := len(ciphertext) - macSize
	//Split Ciphertext from the MAC
	ct, mac := ciphertext[:l], ciphertext[l:]

	// verify mac
	if !poly1305Verify(ct, nonce, &k.MACKey, mac) {
		return nil, ErrUnauthenticated
	}

	ret, out := sliceForAppend(dst, len(ct))

	c, err := aes.NewCipher(k.EncryptionKey[:])
	if err != nil {
		panic(fmt.Sprintf("unable to create cipher: %v", err))
	}
	e := cipher.NewCTR(c, nonce)
	e.XORKeyStream(out, ct)

	return ret, nil
}

/*
Safe validates if data is safe
*/
func (k *Key) Safe(nonce, ciphertext []byte) (bool, error) {
	if !k.Valid() {
		return false, errors.New("invalid key")
	}

	// check parameters
	if len(nonce) != ivSize {
		panic("incorrect nonce length")
	}

	if !validNonce(nonce) {
		return false, errors.New("nonce is invalid")
	}

	// check for plausible length
	if len(ciphertext) < k.Overhead() {
		return false, errors.New("trying to decrypt invalid data: cipher text too small")
	}

	//Size is length of ciphertext minus macSize
	l := len(ciphertext) - macSize
	//Split Ciphertext from the MAC
	ct, mac := ciphertext[:l], ciphertext[l:]

	// verify mac
	if !poly1305Verify(ct, nonce, &k.MACKey, mac) {
		return false, ErrUnauthenticated
	}
	return true, nil
}
