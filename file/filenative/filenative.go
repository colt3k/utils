package filenative

import (
	"os"
	"path/filepath"

	log "github.com/colt3k/nglog/ng"
	"github.com/colt3k/utils/encode/encodeenum"
	"github.com/colt3k/utils/hash"
	"github.com/colt3k/utils/stringut"

	"github.com/colt3k/utils/file"
)

// BasicFile struct
type BasicFile struct {
	fileStr string
	name    string
	size    int64
	file    *os.File
	exists  bool
	hash    []byte
	meta    file.Meta
	ftype   string
}

// NewFile create New nativefile type
func NewFile(file string) file.File {

	t := new(BasicFile)
	if !filepath.IsAbs(file) {
		tmppath, _ := filepath.Abs(file)
		file = tmppath
	}
	t.fileStr = file

	return t
}

// Available does this file exist?
func (f *BasicFile) Available() bool {
	//log.Logln(log.DEBUG, "File:", f.fileStr)

	if _, err := os.Stat(f.fileStr); !os.IsNotExist(err) {
		//log.Logln(log.DEBUG, "File Exists:true")
		return true
	}
	//log.Logln(log.DEBUG, "File Exists:false")
	return false
}

// Path show path to file
func (f *BasicFile) Path() string {
	return f.fileStr
}

/*
Hash creates a hash from the file
*/
func (f *BasicFile) Hash(hasher hash.Hasher, rebuild bool) []byte {

	if len(f.hash) > 0 && !rebuild {
		return f.hash
	}

	if fo, err := os.Open(f.fileStr); err == nil {
		defer fo.Close()
		f.hash = hasher.File(fo)

		return f.hash
	}
	return nil
}

//HRByteCount return a string to represent a size
func (f *BasicFile) HRByteCount(si bool) string {
	return stringut.HRByteCount(f.Size(), si)
}

// SetMeta set meta data on object
func (f *BasicFile) SetMeta(m file.Meta) {
	f.meta = m
}

// Name pull off name from file path
func (f *BasicFile) Name() string {

	if len(f.name) > 0 {
		return f.name
	}
	f.name = filepath.Base(f.fileStr)
	return f.name
}

// Decode file data
func (f *BasicFile) Decode(enctype encodeenum.Encoding) []byte {
	log.Println("not currently implemented")
	return nil
}

// Encode file data
func (f *BasicFile) Encode(enctype encodeenum.Encoding) string {
	log.Println("not currently implemented")
	return ""
}

// PathOnly show only the path no name
func (f *BasicFile) PathOnly() string {
	return filepath.Dir(f.fileStr)
}

// Dir is this a directory
func (f *BasicFile) Dir() bool {
	src, err := os.Stat(f.fileStr)
	if err != nil {
		panic(err)
	}
	return src.IsDir()
}

// Type what is the file type
func (f *BasicFile) Type() string {
	if len(f.ftype) > 0 {
		return f.ftype
	}
	tmp := filepath.Ext(f.Path())
	log.Logln(log.DEBUG, "file type:", tmp)
	if len(tmp) > 0 {
		f.ftype = tmp[1:]
		return f.ftype
	}
	return "unknown"
}

// Size get size of file
func (f BasicFile) Size() int64 {
	if f.size > 0 {
		return f.size
	}
	if fi, err := os.Stat(f.fileStr); err == nil {
		f.size = fi.Size()
		return f.size
	}

	return -1
}

// Delete file
func (f *BasicFile) Delete() bool {
	err := os.Remove(f.fileStr)
	if err == nil {
		return true
	}
	return false
}

// DeleteAll contents
func (f *BasicFile) DeleteAll() bool {
	err := os.RemoveAll(f.fileStr)
	if err == nil {
		return true
	}
	return false
}
