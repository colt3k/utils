package filesize

/*
SizeTypes Enum
*/
type SizeTypes int

const (
	Bytes SizeTypes = 1 + iota
	Kilo
	Mega
	Giga
	Tera
)

var sizeTypes = [...]string{
	"byte(s)", "kb", "mb", "gb", "tb",
}

// Convert from one type size to another
func (sizeType SizeTypes) Convert(val int, origType SizeTypes, si bool) float64 {
	var unit float64 = 1024
	if si {
		unit = 1000
	}
	switch sizeTypes[sizeType-1] {
	case "byte(s)":
		switch sizeTypes[origType-1] {
		case "kb":
			tmp := float64(val) * pow(1, unit)
			return tmp
		case "mb":
			tmp := float64(val) * pow(2, unit)
			return tmp
		case "gb":
			tmp := float64(val) * pow(3, unit)
			return tmp
		case "tb":
			return float64(val) * pow(4, unit)
		}
	case "kb":
		switch sizeTypes[origType-1] {
		case "byte(s)":
			return float64(val) / pow(1, unit)
		case "mb":
			return float64(val) * pow(1, unit)
		case "gb":
			return float64(val) * pow(2, unit)
		case "tb":
			return float64(val) * pow(3, unit)
		}
	case "mb":
		switch sizeTypes[origType-1] {
		case "byte(s)":
			tmp := float64(val) / pow(2, unit)
			return tmp
		case "kb":
			return float64(val) / pow(1, unit)
		case "gb":
			return float64(val) * pow(1, unit)
		case "tb":
			return float64(val) * pow(2, unit)
		}
	case "gb":
		switch sizeTypes[origType-1] {
		case "byte(s)":
			tmp := float64(val) / pow(3, unit)
			return tmp
		case "kb":
			return float64(val) / pow(2, unit)
		case "mb":
			return float64(val) / pow(1, unit)
		case "tb":
			return float64(val) * pow(1, unit)
		}
	case "tb":
		switch sizeTypes[origType-1] {
		case "byte(s)":
			tmp := float64(val) / pow(4, unit)
			return tmp
		case "kb":
			return float64(val) / pow(3, unit)
		case "mb":
			return float64(val) / pow(2, unit)
		case "gb":
			return float64(val) / pow(1, unit)
		}
	}
	return float64(val)
}
func (sizeType SizeTypes) String() string {
	return sizeTypes[sizeType-1]
}

func pow(x int, y float64) float64 {
	var last float64 = 1
	for i := 1; i <= x; i++ {
		last = last * y
	}
	return last
}

/*
Types pulls full list as []string
*/
func (sizeType SizeTypes) Types() []string {
	return sizeTypes[:]
}
