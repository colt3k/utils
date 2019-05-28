package pbkdf2

import (
	"testing"

	log "github.com/colt3k/nglog/ng"
	"github.com/colt3k/utils/encode"
	"github.com/colt3k/utils/encode/encodeenum"
)

func init() {
	log.SetLevel(log.DEBUG)
}

func TestPBKDF2_Generate(t *testing.T) {
	encodedSalt := "qhebZTd7PVqGBCH0rTyl0w=="
	salt := encode.Decode([]byte(encodedSalt), encodeenum.B64STD)
	p := New([]byte("mysupersecretpassword&^%$123"), salt, 0,0)
	key := p.Generate()
	encodedKey := encode.Encode(key, encodeenum.B64STD)
	log.Println(encodedKey)

	if encodedKey != "U9luRG7mGPxQcH3BGOhWfT/amf1glnSKcWitjxPUHvE=" {
		t.FailNow()
	}
}