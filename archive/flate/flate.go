package flate

import (
	"compress/flate"
	"fmt"
	"io"
	"path/filepath"

	"github.com/colt3k/utils/archive"
)

type Flate struct {
	compressionLevel int
}

var FLATE = New()

func New() *Flate {
	return &Flate{
		compressionLevel:flate.DefaultCompression,
	}
}
func (f *Flate) Compress(in io.Reader, out io.Writer) error {
	w, err := flate.NewWriter(out, f.compressionLevel)
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = io.Copy(w, in)
	return err
}
func (f *Flate) Decompress(in io.Reader, out io.Writer) error {
	r := flate.NewReader(in)
	defer r.Close()
	_, err := io.Copy(out, r)
	return err
}
func (f *Flate) Ext(fileName string) error {
	if filepath.Ext(fileName) != ".ft" {
		return fmt.Errorf("must have .ft extension")
	}
	return nil
}

// compile time check to ensure uses interface
var (
	_ = archive.Compressor(new(Flate))
)
