package filesize

import (
	"fmt"
	"testing"

	"github.com/colt3k/utils/mathut"
)

const (
	si = false
)

func assertEqual(t *testing.T, a interface{}, b interface{}, message, typ, toType string) {

	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v for %s to %s", a, b, typ, toType)
	}
	t.Fatal(message)
}

var cases = []struct {
	startSize    SizeTypes
	endSize      SizeTypes
	size         int
	matchSiFalse string
	matchSiTrue  string
	msg          string
}{
	{Bytes, Kilo, 1024, "1", "1.024", ""},
	{Bytes, Mega, 1024, "0.0009765625", "0.001024", ""},
	{Bytes, Giga, 1024, "0.00000095367431640625", "0.000001024", ""},
	{Bytes, Tera, 1024, "0.0000000009313225746154785", "0.000000001024", ""},

	{Kilo, Bytes, 1024, "1048576", "1024000", ""},
	{Kilo, Mega, 1024, "1", "1.024", ""},
	{Kilo, Giga, 1024, "0.0009765625", "0.001024", ""},
	{Kilo, Tera, 1024, "0.00000095367431640625", "0.000001024", ""},

	{Mega, Bytes, 1024, "1073741824", "1024000000", ""},
	{Mega, Kilo, 1024, "1048576", "1024000", ""},
	{Mega, Giga, 1024, "1", "1.024", ""},
	{Mega, Tera, 1024, "0.0009765625", "0.001024", ""},

	{Giga, Bytes, 1024, "1099511627776", "1024000000000", ""},
	{Giga, Kilo, 1024, "1073741824", "1024000000", ""},
	{Giga, Mega, 1024, "1048576", "1024000", ""},
	{Giga, Tera, 1024, "1", "1.024", ""},

	{Tera, Bytes, 1024, "1125899906842624", "1024000000000000", ""},
	{Tera, Kilo, 1024, "1099511627776", "1024000000000", ""},
	{Tera, Mega, 1024, "1073741824", "1024000000", ""},
	{Tera, Giga, 1024, "1048576", "1024000", ""},
}

func TestSizeTypes_Convert(t *testing.T) {

	for _, d := range cases {
		val := SizeTypes(d.endSize).Convert(d.size, d.startSize, false)
		assertEqual(t, mathut.FmtFloat(val), d.matchSiFalse, d.msg, d.startSize.String(), d.endSize.String())

		val = SizeTypes(d.endSize).Convert(d.size, d.startSize, true)
		assertEqual(t, mathut.FmtFloat(val), d.matchSiTrue, d.msg, d.startSize.String(), d.endSize.String())
	}

}
