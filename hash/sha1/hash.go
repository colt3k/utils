package sha1

import (
	"crypto"
	"crypto/sha1"
	"os"

	"github.com/colt3k/utils/hash"
	"github.com/colt3k/utils/hash/hashenum"

	log "github.com/colt3k/nglog/ng"
)

type sha1Hash struct {
	hash   []byte
	format hashenum.HashEnum
}

func NewHash(opts ...HashOption) hash.Hasher {
	h := new(sha1Hash)
	for _, opt := range opts {
		opt(h)
	}

	if h.format != hashenum.SHA1 {
		log.Logln(log.ERROR, "invalid hash for type", h.format.String())
	}

	return h
}
func (h *sha1Hash) Format() hashenum.HashEnum {
	log.Logln(log.DEBUG, "format")
	return h.format
}

func (h *sha1Hash) String(data string) []byte {

	switch h.format {
	case hashenum.SHA1:
		fallthrough
	default:
		if crypto.SHA1.Available() {
			hashType := sha1.New()
			return hash.String(data, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	}

	return nil
}

//HashFile has the file contents
func (h *sha1Hash) File(f *os.File) []byte {

	switch h.format {
	case hashenum.SHA1:
		fallthrough
	default:
		if crypto.SHA1.Available() {
			hashType := sha1.New()
			return hash.File(f, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	}

	return nil
}
