package blake2

import (
	"github.com/colt3k/utils/hash/hashenum"
)

type HashOption func(h *blake2Hash)

func Format(f hashenum.HashEnum) HashOption {
	return func(h *blake2Hash) {
		h.format = f
	}
}
