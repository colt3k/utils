// +build windows

package fileperm

import (
	"os"
	"runtime"
	"syscall"

	log "github.com/colt3k/nglog/ng"
)

//IsHidden determine if file is hidden on Windows
func IsHidden(file os.File) (bool, error) {
	if runtime.GOOS == "windows" {

		pointer, err := syscall.UTF16PtrFromString(file.Name())
		if err != nil {
			return false, err
		}
		attributes, err := syscall.GetFileAttributes(pointer)
		if err != nil {
			return false, err
		}
		return attributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0, nil
	} else {
		log.Logln(log.FATAL, "Unable to check if file is hidden under this OS")
	}
	return false, nil
}
