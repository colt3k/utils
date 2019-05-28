// Code generated by go generate
// This file was generated by robots at 2018-05-01 19:35:01.447367847 +0000 UTC
package hashenum

// HashEnum is an optional HashEnum
type HashEnum int

const (
	MD4 HashEnum = 1 + iota
	MD5
	RIPEMD160
	SHA1
	SHA224
	SHA256
	SHA384
	SHA512
	SHA3_224
	SHA3_256
	SHA3_384
	SHA3_512
	SHA512_224
	SHA512_256
	BLAKE2s_256
	BLAKE2b_256
	BLAKE2b_384
	BLAKE2b_512
)

var hashenum = [...]string{
	"MD4", "MD5", "RIPEMD160", "SHA1", "SHA224", "SHA256", "SHA384", "SHA512", "SHA3_224", "SHA3_256", "SHA3_384", "SHA3_512", "SHA512_224", "SHA512_256", "BLAKE2s_256", "BLAKE2b_256", "BLAKE2b_384", "BLAKE2b_512",
}

func (h HashEnum) String() string {
	return hashenum[h-1]
}

/*
Types pulls full list as []string
*/
func (h HashEnum) Types() []string {
	return hashenum[:]
}
