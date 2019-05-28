package md

import (
	"crypto"
	"crypto/md5"
	"os"

	log "github.com/colt3k/nglog/ng"

	"github.com/colt3k/utils/hash"
	"github.com/colt3k/utils/hash/hashenum"
)

type mdHash struct {
	hash   []byte
	format hashenum.HashEnum
}

func NewHash(opts ...HashOption) hash.Hasher {

	h := new(mdHash)
	for _, opt := range opts {
		opt(h)
	}

	if h.format > hashenum.RIPEMD160 {
		log.Logln(log.ERROR, "invalid hash for type", h.format.String())
	}

	return h
}

func (h *mdHash) Format() hashenum.HashEnum {
	log.Logln(log.DEBUG, "format")
	return h.format
}
func (h *mdHash) String(data string) []byte {
	switch h.format {
	case hashenum.RIPEMD160:
		if crypto.RIPEMD160.Available() {
			hashType := crypto.RIPEMD160.New()
			return hash.String(data, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available", 1.0, true, 'R', byte(2))
	case hashenum.MD5:
		fallthrough
	default:
		if crypto.MD5.Available() {
			hashType := md5.New()
			return hash.String(data, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	}

	return nil
}

//HashFile has the file contents
func (h *mdHash) File(f *os.File) []byte {

	switch h.format {
	case hashenum.RIPEMD160:
		if crypto.RIPEMD160.Available() {
			hashType := crypto.RIPEMD160.New()
			return hash.File(f, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	case hashenum.MD5:
		fallthrough
	default:
		if crypto.MD5.Available() {
			hashType := md5.New()
			return hash.File(f, hashType)
		}
		log.Logln(log.ERROR, h.format.String(), "not available")
	}

	return nil
}
