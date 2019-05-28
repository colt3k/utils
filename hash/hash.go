package hash

import (
	"hash"
	"io"
	"os"

	"github.com/colt3k/utils/hash/hashenum"
)

//go:generate enumeration -pkg hashenum -type HashEnum -list MD4,MD5,RIPEMD160,SHA1,SHA224,SHA256,SHA384,SHA512,SHA3_224,SHA3_256,SHA3_384,SHA3_512,SHA512_224,SHA512_256,BLAKE2s_256,BLAKE2b_256,BLAKE2b_384,BLAKE2b_512

type Hasher interface {
	String(data string) []byte
	File(f *os.File) []byte
	Format() hashenum.HashEnum
}

func String(data string, hash hash.Hash) []byte {
	io.WriteString(hash, data)
	val := hash.Sum(nil)
	return val
}

func File(file *os.File, hash hash.Hash) []byte {
	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		panic(err)
	}
	//Get the 16 bytes hash
	val := hash.Sum(nil)
	return val
}
