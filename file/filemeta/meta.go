package filemeta

import (
	"os"
	"runtime"
	"time"

	log "github.com/colt3k/nglog/ng"

	"github.com/colt3k/utils/file"
	"github.com/colt3k/utils/file/fileperm"
	"github.com/colt3k/utils/file/mimetypes"
)

type BaseFileMeta struct {
	file     file.File
	ftype    string `json:"file_ext"`
	mimetype string
	perms    fileperm.PermissionBits
	mtime    time.Time
	atime    time.Time
	ctime    time.Time
	lmod     int64 `json:"lastmod"`
}

func New(fileObj file.File) file.Meta {
	t := new(BaseFileMeta)
	t.file = fileObj
	return t
}

func (b BaseFileMeta) Hidden() bool {
	if runtime.GOOS != "windows" {
		//if fo, err := b.file.Open(); err == nil {
		//	//Tell the program to call the following function when the current function returns
		//	defer fo.Close()
		//
		//	// unix/linux file or directory that starts with . is hidden
		//	if fo.Name()[0:1] == "." {
		//		return true
		//	}
		//}
	} else {
		log.Logln(log.FATAL, "Unable to check if file is hidden under this OS")
	}
	return false
}

func (b BaseFileMeta) NamedPipe() bool {
	src, err := os.Lstat(b.file.Path())
	if err != nil {
		panic(err)
	}
	return src.Mode()&os.ModeNamedPipe != 0
}
func (b BaseFileMeta) Regular() bool {
	src, err := os.Stat(b.file.Path())
	if err != nil {
		panic(err)
	}
	return src.Mode().IsRegular()
}
func (b BaseFileMeta) SymLink() bool {
	src, err := os.Lstat(b.file.Path())
	if err != nil {
		panic(err)
	}
	return src.Mode()&os.ModeSymlink != 0
}

//LastMod last modification date/time on file
func (b BaseFileMeta) LastMod() int64 {
	_, mTime, _, _ := b.Times()
	//log.Print("MTime as Millis:",(mTime.UnixNano()/(int64(time.Millisecond))))
	b.lmod = mTime.UnixNano() / (int64(time.Millisecond))

	return b.lmod
}
func (b BaseFileMeta) MetaType() string {
	if len(b.mimetype) > 0 {
		return b.mimetype
	}
	b.mimetype = mimetypes.Find(b.ftype)

	return b.mimetype
}

//Permissions retrieve permissions on file
func (b BaseFileMeta) Permissions() fileperm.PermissionBits {

	if len(b.perms.String()) > 0 {
		return b.perms
	}
	fi, err := os.Stat(b.file.Path())
	if err != nil {
		panic(err)
	}
	fm := fi.Mode()
	return fileperm.PermissionBits(fm.Perm())

}
