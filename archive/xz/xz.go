package xz

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/colt3k/utils/archive"
	"github.com/ulikunitz/xz"
	fastxz "github.com/xi2/xz"
)

type Xz struct{}

var XZ = NewXz()

func NewXz() *Xz {
	return new(Xz)
}

// Compress reads in, compresses it, and writes it to out.
func (x *Xz) Compress(in io.Reader, out io.Writer) error {
	w, err := xz.NewWriter(out)
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = io.Copy(w, in)
	return err
}

// Decompress reads in, decompresses it, and writes it to out.
func (x *Xz) Decompress(in io.Reader, out io.Writer) error {
	r, err := fastxz.NewReader(in, 0)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, r)
	return err
}

// CheckExt ensures the file extension matches the format.
func (x *Xz) CheckExt(filename string) error {
	if filepath.Ext(filename) != ".xz" {
		return fmt.Errorf("must have a .xz extension")
	}
	return nil
}

// compile time check to ensure uses interface
var (
	_ = archive.Compressor(new(Xz))
)
