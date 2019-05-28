package archive

import "io"

type Compressor interface {
	Compress(in io.Reader, out io.Writer) error
	Decompress(in io.Reader, out io.Writer) error
}
