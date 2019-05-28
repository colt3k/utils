package lzw

import (
	"compress/gzip"
	"compress/lzw"
	"fmt"
	"io"
	"path/filepath"

	"github.com/colt3k/utils/archive"
)

type Lzw struct {
	litWidth int
	order    lzw.Order
}

var LZW = New()

func New() *Lzw {
	return &Lzw{
		litWidth:8,
		order:lzw.LSB,
	}
}
func (l *Lzw) Compress(in io.Reader, out io.Writer) error {
	w := lzw.NewWriter(out, l.order, l.litWidth)

	defer w.Close()
	_, err := io.Copy(w, in)
	return err
}
func (l *Lzw) Decompress(in io.Reader, out io.Writer) error {
	r, err := gzip.NewReader(in)
	if err != nil {
		return err
	}
	defer r.Close()
	_, err = io.Copy(out, r)
	return err
}
func (l *Lzw) Ext(fileName string) error {
	if filepath.Ext(fileName) != ".lzw" {
		return fmt.Errorf("must have .lzw extension")
	}
	return nil
}

// compile time check to ensure uses interface
var (
	_ = archive.Compressor(new(Lzw))
)
