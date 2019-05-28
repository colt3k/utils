package sha2

import (
	"crypto"
	"crypto/sha256"
	"os"

	log "github.com/colt3k/nglog/ng"

	"github.com/colt3k/utils/hash"
	"github.com/colt3k/utils/hash/hashenum"
)

type sha2Hash struct {
	hash   []byte
	format hashenum.HashEnum
}

func NewHash(opts ...HashOption) hash.Hasher {
	h := new(sha2Hash)
	for _, opt := range opts {
		opt(h)
	}

	if h.format < hashenum.SHA224 || h.format > hashenum.SHA256 {
		log.Logln(log.ERROR, "invalid hash for type", h.format.String())
	}
	return h
}
func (h *sha2Hash) Format() hashenum.HashEnum {
	log.Logln(log.DEBUG, "format")
	return h.format
}
func (h *sha2Hash) String(data string) []byte {

	switch h.format {
	case hashenum.SHA224:
		if crypto.SHA224.Available() {
			hashType := sha256.New224()
			return hash.String(data, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	case hashenum.SHA256:
		if crypto.SHA256.Available() {
			hashType := sha256.New()
			return hash.String(data, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	}

	return nil
}

//HashFile has the file contents
func (h *sha2Hash) File(f *os.File) []byte {
	switch h.format {
	case hashenum.SHA224:
		if crypto.SHA224.Available() {
			hashType := sha256.New224()
			return hash.File(f, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	case hashenum.SHA256:
		if crypto.SHA256.Available() {
			hashType := sha256.New()
			return hash.File(f, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	}

	return nil
}
