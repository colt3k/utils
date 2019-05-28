package encode

//go:generate enumeration -pkg enodeenum -type Encoding -list B64STD,B64URL,Hex -hrtypes b64standard,b64url,hex

import (
	"encoding/base64"
	"encoding/hex"
	"strings"

	log "github.com/colt3k/nglog/ng"
	"github.com/colt3k/utils/encode/encodeenum"
)

//B64DecodeStdSanitized decode base 64 and sanitize (backtick, double quotes)
func B64DecodeStdSanitized(data string) []byte {
	sanitized := strings.Replace(data, `\`, "", -1)
	sani := []byte(sanitized)
	return Decode(sani, encodeenum.B64STD)
}

//Encode process encoding for the type passed
func Encode(data []byte, enctype encodeenum.Encoding) string {

	switch enctype {
	case encodeenum.Hex:
		tmp := hex.EncodeToString(data)
		return tmp
	case encodeenum.B64STD:
		tmp := base64.StdEncoding.EncodeToString(data)
		return tmp
	case encodeenum.B64URL:
		tmp := base64.URLEncoding.EncodeToString(data)
		return tmp
	}

	log.Logln(log.ERROR, "Invalid Encoding Type Passed, value returned as string")

	tmp := string(data)
	return tmp
}

//Decode process decoding for the type passed in
func Decode(data []byte, enctype encodeenum.Encoding) []byte {

	switch enctype {
	case encodeenum.Hex:
		hexdata := make([]byte, hex.DecodedLen(len(data)))
		l, _ := hex.Decode(hexdata, data)
		tmp := hexdata[:l]
		return tmp
	case encodeenum.B64STD:
		base64data := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
		n, _ := base64.StdEncoding.Decode(base64data, data)
		tmp := base64data[:n]
		return tmp
	case encodeenum.B64URL:
		base64data := make([]byte, base64.URLEncoding.DecodedLen(len(data)))
		n, _ := base64.URLEncoding.Decode(base64data, data)
		tmp := base64data[:n]
		return tmp
	}

	log.Logln(log.ERROR, "Invalid Decoding Type Passed, value returned as nil")

	return nil
}
