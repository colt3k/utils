package sha3

import (
	"crypto"
	"os"

	"golang.org/x/crypto/sha3"

	log "github.com/colt3k/nglog/ng"

	"github.com/colt3k/utils/hash"
	"github.com/colt3k/utils/hash/hashenum"
)

type sha3Hash struct {
	hash   []byte
	format hashenum.HashEnum
}

func NewHash(opts ...HashOption) hash.Hasher {
	h := new(sha3Hash)
	for _, opt := range opts {
		opt(h)
	}

	if h.format < hashenum.SHA3_224 || h.format > hashenum.SHA3_512 {
		log.Logln(log.ERROR, "invalid hash for type", h.format.String())
	}

	return h
}
func (h *sha3Hash) Format() hashenum.HashEnum {
	log.Logln(log.DEBUG, "format")
	return h.format
}
func (h *sha3Hash) String(data string) []byte {

	switch h.format {
	case hashenum.SHA3_224:
		if crypto.SHA3_224.Available() {
			hashType := sha3.New224()
			return hash.String(data, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	case hashenum.SHA3_256:
		if crypto.SHA3_256.Available() {
			hashType := sha3.New256()
			return hash.String(data, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	case hashenum.SHA3_384:
		if crypto.SHA3_384.Available() {
			hashType := sha3.New384()
			return hash.String(data, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	case hashenum.SHA3_512:
		if crypto.SHA3_512.Available() {
			hashType := sha3.New512()
			return hash.String(data, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	}

	return nil
}

//HashFile has the file contents
func (h *sha3Hash) File(f *os.File) []byte {

	switch h.format {
	case hashenum.SHA3_224:
		if crypto.SHA3_224.Available() {
			hashType := sha3.New224()
			return hash.File(f, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	case hashenum.SHA3_256:
		if crypto.SHA3_256.Available() {
			hashType := sha3.New256()
			return hash.File(f, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	case hashenum.SHA3_384:
		if crypto.SHA3_384.Available() {
			hashType := sha3.New384()
			return hash.File(f, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	case hashenum.SHA3_512:
		if crypto.SHA3_512.Available() {
			hashType := sha3.New512()
			return hash.File(f, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	}

	return nil
}
