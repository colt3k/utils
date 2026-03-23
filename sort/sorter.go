// Package utils in the sort module provides a case-insensitive sorter for
// store.FileStore slices.
package utils

import (
	"unicode"

	"github.com/colt3k/utils/store"
)

type Stores []store.FileStore

// Len returns the number of stored file records.
func (s Stores) Len() int {
	return len(s)
}

// Swap swaps two items in the slice.
func (s Stores) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less compares file names case-insensitively, using original rune case as a
// tie-breaker.
func (s Stores) Less(i, j int) bool {
	iRunes := []rune(s[i].Name)
	jRunes := []rune(s[j].Name)

	max := len(iRunes)
	if max > len(jRunes) {
		max = len(jRunes)
	}

	for idx := 0; idx < max; idx++ {
		ir := iRunes[idx]
		jr := jRunes[idx]

		lir := unicode.ToLower(ir)
		ljr := unicode.ToLower(jr)

		if lir != ljr {
			return lir < ljr
		}

		// the lowercase runes are the same, so compare the original
		if ir != jr {
			return ir < jr
		}
	}

	return false
}
