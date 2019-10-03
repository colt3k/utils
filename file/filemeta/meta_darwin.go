package filemeta

import (
	"os"
	"syscall"
	"time"
)


func (b BaseFileMeta) Times() (atime, mtime, ctime time.Time, err error) {
	fi, err := os.Stat(b.file.Path())
	if err != nil {
		return
	}
	b.mtime = fi.ModTime()
	stat := fi.Sys().(*syscall.Stat_t)
	b.atime = time.Unix(int64(stat.Atimespec.Sec), int64(stat.Atimespec.Nsec))
	b.ctime = time.Unix(int64(stat.Ctimespec.Sec), int64(stat.Ctimespec.Nsec))
	return b.atime, b.mtime, b.ctime, nil
}
