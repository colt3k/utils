package sha2

import (
	"github.com/colt3k/utils/hash/hashenum"
)

type HashOption func(h *sha2Hash)

func Format(f hashenum.HashEnum) HashOption {
	return func(h *sha2Hash) {
		h.format = f
	}
}
