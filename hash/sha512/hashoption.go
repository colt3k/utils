package sha512

import (
	"github.com/colt3k/utils/hash/hashenum"
)

type HashOption func(h *sha512Hash)

func Format(f hashenum.HashEnum) HashOption {
	return func(h *sha512Hash) {
		h.format = f
	}
}
