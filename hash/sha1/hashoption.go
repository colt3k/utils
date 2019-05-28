package sha1

import (
	"github.com/colt3k/utils/hash/hashenum"
)

type HashOption func(h *sha1Hash)

func Format(f hashenum.HashEnum) HashOption {
	return func(h *sha1Hash) {
		h.format = f
	}
}
