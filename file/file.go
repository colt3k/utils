package file

import (
	"os"
	"time"

	"github.com/colt3k/utils/encode/encodeenum"
	"github.com/colt3k/utils/hash"

	"github.com/colt3k/utils/file/fileperm"
)

type Meta interface {
	Hidden() bool
	NamedPipe() bool
	Regular() bool
	SymLink() bool
	LastMod() int64
	Permissions() fileperm.PermissionBits
	Times() (time.Time, time.Time, time.Time, error)
}
type Exif interface {
	Exif(fileName string)
	ReadLatLongData(f *os.File)
	ReadALLDataAsJSON(f *os.File)
}

type Image interface {
	Dimensions(fileName string) (int, int)
}
type File interface {
	Available() bool
	Delete() bool
	DeleteAll() bool
	Dir() bool
	Type() string
	Path() string
	PathOnly() string
	Name() string
	Size() int64
	Hash(hasher hash.Hasher, rebuild bool) []byte
	HRByteCount(si bool) string
	SetMeta(m Meta)
	Decode(enctype encodeenum.Encoding) []byte
	Encode(enctype encodeenum.Encoding) string
}
