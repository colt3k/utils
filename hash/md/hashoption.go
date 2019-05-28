package md

import (
	"github.com/colt3k/utils/hash/hashenum"
)

type HashOption func(h *mdHash)

func Format(f hashenum.HashEnum) HashOption {
	return func(h *mdHash) {

		h.format = f

	}
}
