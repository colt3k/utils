package blake2

import (
	"crypto"
	"os"

	log "github.com/colt3k/nglog/ng"

	"github.com/colt3k/utils/hash"
	"github.com/colt3k/utils/hash/hashenum"
)

type blake2Hash struct {
	hash   []byte
	format hashenum.HashEnum
}

func NewHash(opts ...HashOption) hash.Hasher {
	h := new(blake2Hash)
	for _, opt := range opts {
		opt(h)
	}
	// Default logger
	if h.format < hashenum.BLAKE2s_256 || h.format > hashenum.BLAKE2b_512 {
		log.Logln(log.ERROR, "invalid hash for type", h.format.String())
	}

	return h
}

func (h *blake2Hash) Format() hashenum.HashEnum {
	log.Logln(log.DEBUG, "format")
	return h.format
}
func (h *blake2Hash) String(data string) []byte {

	switch h.format {
	case hashenum.BLAKE2s_256:
		if crypto.BLAKE2s_256.Available() {
			hashType := crypto.BLAKE2s_256.New()
			return hash.String(data, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	case hashenum.BLAKE2b_256:
		if crypto.BLAKE2b_256.Available() {
			hashType := crypto.BLAKE2b_256.New()
			return hash.String(data, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	case hashenum.BLAKE2b_384:
		if crypto.BLAKE2b_384.Available() {
			hashType := crypto.BLAKE2b_384.New()
			return hash.String(data, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	case hashenum.BLAKE2b_512:
		if crypto.BLAKE2b_512.Available() {
			hashType := crypto.BLAKE2b_512.New()
			return hash.String(data, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	}

	return nil
}

//HashFile has the file contents
func (h *blake2Hash) File(f *os.File) []byte {

	switch h.format {
	case hashenum.BLAKE2s_256:
		if crypto.BLAKE2s_256.Available() {
			hashType := crypto.BLAKE2s_256.New()
			return hash.File(f, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	case hashenum.BLAKE2b_256:
		if crypto.BLAKE2b_256.Available() {
			hashType := crypto.BLAKE2b_256.New()
			return hash.File(f, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	case hashenum.BLAKE2b_384:
		if crypto.BLAKE2b_384.Available() {
			hashType := crypto.BLAKE2b_384.New()
			return hash.File(f, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	case hashenum.BLAKE2b_512:
		if crypto.BLAKE2b_512.Available() {
			hashType := crypto.BLAKE2b_512.New()
			return hash.File(f, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	}

	return nil
}
