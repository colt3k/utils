package gz

import (
	"compress/gzip"
	"fmt"
	"io"
	"path/filepath"

	"github.com/colt3k/utils/archive"
)

type GZip struct {
	compressionLevel int
}

var Gz = New()

func New() *GZip {
	return &GZip{
		compressionLevel:gzip.DefaultCompression,
	}
}
func (g *GZip) Compress(in io.Reader, out io.Writer) error {
	w, err := gzip.NewWriterLevel(out, g.compressionLevel)
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = io.Copy(w, in)
	return err
}
func (g *GZip) Decompress(in io.Reader, out io.Writer) error {
	r, err := gzip.NewReader(in)
	if err != nil {
		return err
	}
	defer r.Close()
	_, err = io.Copy(out, r)
	return err
}
func (g *GZip) Ext(fileName string) error {
	if filepath.Ext(fileName) != ".gz" {
		return fmt.Errorf("must have .gz extension")
	}
	return nil
}

// compile time check to ensure uses interface
var (
	_ = archive.Compressor(new(GZip))
)