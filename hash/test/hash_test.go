package test

import (
	log "github.com/colt3k/nglog/ng"
	"github.com/colt3k/utils/encode"
	"github.com/colt3k/utils/encode/encodeenum"
	"github.com/colt3k/utils/hash"
	"github.com/colt3k/utils/hash/blake2"
	"github.com/colt3k/utils/hash/hashenum"
	"github.com/colt3k/utils/hash/md"
	"github.com/colt3k/utils/hash/sha1"
	"github.com/colt3k/utils/hash/sha2"
	"github.com/colt3k/utils/hash/sha3"
	"github.com/colt3k/utils/hash/sha512"
)

/*
MD5
SHA1
SHA256
*/

var (
	plain = "this is my example text"
)

func ExampleOne() {

	mdFmt := md.NewHash(md.Format(hashenum.MD5))
	show(mdFmt, mdFmt.String(plain))
	mdFmt = md.NewHash(md.Format(hashenum.RIPEMD160))
	show(mdFmt, mdFmt.String(plain))

	sha1Fmt := sha1.NewHash(sha1.Format(hashenum.SHA1))
	show(sha1Fmt, sha1Fmt.String(plain))

	sha2Fmt := sha2.NewHash(sha2.Format(hashenum.SHA256))
	show(sha2Fmt, sha2Fmt.String(plain))

	sha3Fmt := sha3.NewHash(sha3.Format(hashenum.SHA3_224))
	show(sha3Fmt, sha3Fmt.String(plain))
	sha3Fmt = sha3.NewHash(sha3.Format(hashenum.SHA3_256))
	show(sha3Fmt, sha3Fmt.String(plain))
	sha3Fmt = sha3.NewHash(sha3.Format(hashenum.SHA3_512))
	show(sha3Fmt, sha3Fmt.String(plain))

	sha512Fmt := sha512.NewHash(sha512.Format(hashenum.SHA384))
	show(sha512Fmt, sha512Fmt.String(plain))
	sha512Fmt = sha512.NewHash(sha512.Format(hashenum.SHA512))
	show(sha512Fmt, sha512Fmt.String(plain))
	sha512Fmt = sha512.NewHash(sha512.Format(hashenum.SHA512_224))
	show(sha512Fmt, sha512Fmt.String(plain))
	sha512Fmt = sha512.NewHash(sha512.Format(hashenum.SHA512_256))
	show(sha512Fmt, sha512Fmt.String(plain))

	blakeFmt := blake2.NewHash(blake2.Format(hashenum.BLAKE2s_256))
	show(blakeFmt, blakeFmt.String(plain))
	blakeFmt = blake2.NewHash(blake2.Format(hashenum.BLAKE2b_256))
	show(blakeFmt, blakeFmt.String(plain))
	blakeFmt = blake2.NewHash(blake2.Format(hashenum.BLAKE2b_384))
	show(blakeFmt, blakeFmt.String(plain))
	blakeFmt = blake2.NewHash(blake2.Format(hashenum.BLAKE2b_512))
	show(blakeFmt, blakeFmt.String(plain))
	/*
		Output:
		will fail
	*/
}

func show(hasher hash.Hasher, data []byte) {

	if len(hasher.Format().String()) <= 4 {
		log.Println("\tTYPE:", hasher.Format().String(), "\t\t\t|HASHED_DATA:", encode.Encode(data, encodeenum.B64STD), "|")
	} else if len(hasher.Format().String()) <= 8 {
		log.Println("\tTYPE:", hasher.Format().String(), "\t\t|HASHED_DATA:", encode.Encode(data, encodeenum.B64STD), "|")
	} else {
		log.Println("\tTYPE:", hasher.Format().String(), "\t|HASHED_DATA:", encode.Encode(data, encodeenum.B64STD), "|")
	}
	if data == nil {
		log.SetFlags(0)
		log.Logln(log.NONE, "--- previous error ---")
	}
}
