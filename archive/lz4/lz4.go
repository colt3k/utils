package lz4

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/pierrec/lz4"
	"github.com/colt3k/utils/archive"
)

type Lz4 struct {
	compressionLevel int
}

var LZ4 = NewLz4()

func NewLz4() *Lz4 {
	return &Lz4{
		compressionLevel: 9, // https://github.com/lz4/lz4/blob/1b819bfd633ae285df2dfe1b0589e1ec064f2873/lib/lz4hc.h#L48
	}
}

func (lz *Lz4) Compress(in io.Reader, out io.Writer) error {
	w := lz4.NewWriter(out)
	w.Header.CompressionLevel = lz.compressionLevel
	defer w.Close()
	_, err := io.Copy(w, in)
	return err
}

// Decompress reads in, decompresses it, and writes it to out.
func (lz *Lz4) Decompress(in io.Reader, out io.Writer) error {
	r := lz4.NewReader(in)
	_, err := io.Copy(out, r)
	return err
}

// CheckExt ensures the file extension matches the format.
func (lz *Lz4) CheckExt(filename string) error {
	if filepath.Ext(filename) != ".lz4" {
		return fmt.Errorf("must have a .lz4 extension")
	}
	return nil
}

// compile time check to ensure uses interface
var (
	_ = archive.Compressor(new(Lz4))
)