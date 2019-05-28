package mathut

import (
	"math"
	"strconv"
)

func Round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

// The format fmt is one of
// 'b' (-ddddp±ddd, a binary exponent),
// 'e' (-d.dddde±dd, a decimal exponent),
// 'E' (-d.ddddE±dd, a decimal exponent),
// 'f' (-ddd.dddd, no exponent),
// 'g' ('e' for large exponents, 'f' otherwise), or
// 'G' ('E' for large exponents, 'f' otherwise).
func FmtFloat(val float64) string {
	return FmtFloatWithPrecision(val,  -1)
}
func FmtFloatWithPrecision(val float64, precision int) string {
	return strconv.FormatFloat(val, 'f', precision, 64)
}
func FmtFloatExpon(val float64) string {
	return strconv.FormatFloat(val, 'E', -1, 64)
}
// FmtInt parse 32 or 64 cast to int
func FmtInt(val int) string {
	switch IntSize() {
	case 32:
		return strconv.Itoa(val)
	case 64:
		return strconv.FormatInt(int64(val), 10)
	}
	return strconv.Itoa(val)
}
// ParseFloat size 32 or 64
func ParseFloat(val string, sz int) (float64,error) {
	return strconv.ParseFloat(val, sz)
}

func IntSize() int {
	return strconv.IntSize
}
func ParseInt(val string) int64 {
	v, _ := strconv.ParseInt(val, 10, 64)
	return v
}
/* Bitwise operations

	1101 =
	1*2³ + 1*2² + 0*2¹ 	+ 1*2⁰ 	=
	8 	 + 4 	+ 0 	+ 1		=  13

	The ^ operator does a bitwise complement, flips bits from 1 to 0 and 0 to 1
		for example with 3 unsigned bits, ^(101) = 010

	The >> is the right shift operator, a right shift moves all of the bits to the right,
		dropping bits off the right and inserting zeros on the left.
		Example:  3 unsigned bits, 101 >> 2 = 001

	The << is the left shift operator just like the right except that bits shift
		the opposite direction.  Example:	101 << 2 = 100

	Here is how these operators are used in the strconv.IntSize expression

	Expression				32 bit representation		64 bit representation
	uint(0)					00...00 (32 zeros)			0000...0000 (64 zeros)
	^uint(0)				11...11 (32 ones)			1111...1111 (64 ones)
	(^uint(0) >> 63)		00...00 = 0					0000...0001 = 1
	32						00...100000					0000...100000
	32 << (^uint(0) >> 63)	32 << (0)					32 << (1)
							= 100000 << 0				= 100000 << 1
							= 32						= 1000000
														= 64

	In other words
	1. Start with 0
	2. ^ to flip all bits to 1
	3. Right shift (>>) by 63 to only keep a single 1 from 64-bit numbers and zero out 32-bit numbers.
	4. Left shift (<<) 32 by whatever the result is.
	5. This leaves 32 on architectures that use 32-bit integer representations and 64 for 64-bit architectures.


 */