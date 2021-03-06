package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/colt3k/utils/encode"
	"github.com/colt3k/utils/encode/encodeenum"

	"github.com/colt3k/utils/archive/xz"
)

func main() {

	str := "just some regular text to compress and decompress, just some regular text to compress and decompress"
	fmt.Println("Original: ",str)

	var byt bytes.Buffer
	xz.XZ.Compress(strings.NewReader(str), &byt)

	enc := encode.Encode(byt.Bytes(), encodeenum.B64STD)
	fmt.Println("Compressed and Encoded: ",enc)

	dec := encode.Decode([]byte(enc), encodeenum.B64STD)
	fmt.Println("Compressed and Decoded: ",string(dec))

	var b2 bytes.Buffer
	xz.XZ.Decompress(bytes.NewReader(dec), &b2)
	fmt.Println("De-compressed: ", string(b2.String()))
}