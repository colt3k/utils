// Package archive defines the shared compression interface implemented by the
// format-specific packages in this module.
package archive

import "io"

// Compressor describes a stream-based compressor/decompressor pair.
type Compressor interface {
	Compress(in io.Reader, out io.Writer) error
	Decompress(in io.Reader, out io.Writer) error
}
