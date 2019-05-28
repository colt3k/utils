package sha512

import (
	"crypto"
	"crypto/sha512"
	"os"

	log "github.com/colt3k/nglog/ng"

	"github.com/colt3k/utils/hash"
	"github.com/colt3k/utils/hash/hashenum"
)

type sha512Hash struct {
	hash   []byte
	format hashenum.HashEnum
}

func NewHash(opts ...HashOption) hash.Hasher {
	h := new(sha512Hash)
	for _, opt := range opts {
		opt(h)
	}

	if h.format != hashenum.SHA384 && h.format != hashenum.SHA512 && (h.format < hashenum.SHA512_224 ||
		h.format > hashenum.SHA512_256) {
		log.Logln(log.ERROR, "invalid hash for type", h.format.String())
	}

	return h
}
func (h *sha512Hash) Format() hashenum.HashEnum {
	log.Logln(log.DEBUG, "format")
	return h.format
}
func (h *sha512Hash) String(data string) []byte {

	switch h.format {
	case hashenum.SHA384:
		if crypto.SHA384.Available() {
			hashType := sha512.New384()
			return hash.String(data, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	case hashenum.SHA512:
		if crypto.SHA512.Available() {
			hashType := sha512.New()
			return hash.String(data, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	case hashenum.SHA512_224:
		if crypto.SHA512_224.Available() {
			hashType := sha512.New512_224()
			return hash.String(data, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	case hashenum.SHA512_256:
		if crypto.SHA512_256.Available() {
			hashType := sha512.New512_256()
			return hash.String(data, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	}

	return nil
}

//HashFile has the file contents
func (h *sha512Hash) File(f *os.File) []byte {

	switch h.format {
	case hashenum.SHA384:
		if crypto.SHA384.Available() {
			hashType := sha512.New384()
			return hash.File(f, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	case hashenum.SHA512:
		if crypto.SHA512.Available() {
			hashType := sha512.New()
			return hash.File(f, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	case hashenum.SHA512_224:
		if crypto.SHA512_224.Available() {
			hashType := sha512.New512_224()
			return hash.File(f, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	case hashenum.SHA512_256:
		if crypto.SHA512_256.Available() {
			hashType := sha512.New512_256()
			return hash.File(f, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	}

	return nil
}
