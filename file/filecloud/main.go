package filecloud

import (
	"path/filepath"

	log "github.com/colt3k/nglog/ng"
	"github.com/colt3k/utils/encode/encodeenum"
	"github.com/colt3k/utils/hash"

	"github.com/colt3k/utils/file"
	"github.com/colt3k/utils/file/filesize"
)

type CloudFile struct {
	fileStr string
	name    string
	size    int64
	hash    []byte
	meta    file.Meta
	ftype   string
}

func NewFile(file string) file.File {

	t := new(CloudFile)
	t.fileStr = file

	return t
}

func (f *CloudFile) Available() bool {
	return false
}
func (f *CloudFile) Path() string {
	return f.fileStr
}

func (f *CloudFile) SizeAs(szType filesize.SizeTypes, si bool) float64 {

	return 0
}

func (f *CloudFile) Hash(hasher hash.Hasher, rebuild bool) []byte {

	return nil
}

//HRByteCount return a string to represent a size
func (f *CloudFile) HRByteCount(si bool) string {
	return ""
}
func (f *CloudFile) SetMeta(m file.Meta) {
	f.meta = m
}
func (f *CloudFile) Name() string {
	return f.name
}
func (f *CloudFile) Decode(enctype encodeenum.Encoding) []byte {
	log.Println("not currently implemented")
	return nil
}
func (f *CloudFile) Encode(enctype encodeenum.Encoding) string {
	log.Println("not currently implemented")
	return ""
}
func (f *CloudFile) PathOnly() string {
	return filepath.Dir(f.fileStr)
}
func (f *CloudFile) Dir() bool {
	return false
}
func (f CloudFile) Type() string {
	if len(f.ftype) > 0 {
		return f.ftype
	}
	f.ftype = filepath.Ext(f.Path())

	return f.ftype
}
func (f *CloudFile) Delete() bool {

	return false
}
func (f *CloudFile) DeleteAll() bool {

	return false
}
func (f CloudFile) Size() int64 {

	return -1
}
