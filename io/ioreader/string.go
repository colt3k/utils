package ioreader

import (
	"bufio"

	"github.com/iancoleman/orderedmap"

	"github.com/colt3k/utils/io"
)

type String struct {
}

/*
Line returns a single line (without the ending \n)
from the input buffered reader.
An error is returned if` there is an error with the
buffered reader.
*/
func (s *String) ReadLine(br *bufio.Reader) (string, error) {
	return io.ReadLine(br)
}

// Bytes read file into passed byte array
func (s *String) Bytes(buf []byte) {
}

// AsBytes read file into []byte
func (s *String) AsBytes() ([]byte, error) {
	return nil, nil
}

// AsString read file as a string
func (s *String) AsString() (string, error) {
	return "", nil
}

func (s *String) LineScanner(maxBuf, maxScanTokenSize int) (string, error) {
	return "", nil
}

// AsCSVIntoMap read csv file into map[string]string
func (s *String) AsCSVIntoMap() *map[string]string {
	return nil
}

// AsCSVIntoOrderedMap read csv into an ordered map
func (s *String) AsCSVIntoOrderedMap() (*orderedmap.OrderedMap, error) {
	return nil, nil
}
