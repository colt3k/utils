// +build !windows

package fileperm

import (
	"os"
	"runtime"
	"log"
)

func IsHidden(file os.File) bool {
	if runtime.GOOS != "windows" {

		// unix/linux file or directory that starts with . is hidden
		if file.Name()[0:1] == "." {
			return true

		} else {
			return false
		}

	} else {
		log.Fatal("Unable to check if file is hidden under this OS")
	}
	return false
}
