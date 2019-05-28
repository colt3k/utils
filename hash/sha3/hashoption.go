package sha3

import (
	"github.com/colt3k/utils/hash/hashenum"
)

type HashOption func(h *sha3Hash)

func Format(f hashenum.HashEnum) HashOption {
	return func(h *sha3Hash) {
		h.format = f
	}
}
