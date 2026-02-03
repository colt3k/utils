package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/colt3k/utils/archive/lz4"
)

func main() {

	str := "just some regular text to compress and decompress, just some regular text to compress and decompress"
	fmt.Println("Original: ", str)

	var byt bytes.Buffer
	lz4.LZ4.Compress(strings.NewReader(str), &byt)

	enc := base64.StdEncoding.EncodeToString(byt.Bytes())
	fmt.Println("Compressed and Encoded: ", enc)

	base64data := make([]byte, base64.StdEncoding.DecodedLen(len([]byte(enc))))
	n, _ := base64.StdEncoding.Decode(base64data, []byte(enc))
	dec := base64data[:n]
	fmt.Println("Compressed and Decoded: ", string(dec))

	var b2 bytes.Buffer
	lz4.LZ4.Decompress(bytes.NewReader(dec), &b2)
	fmt.Println("De-compressed: ", string(b2.String()))
}
