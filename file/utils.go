package file

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	log "github.com/colt3k/nglog/ng"

	"github.com/colt3k/utils/file/filesize"
)

func PathSeparator() string {
	return string(filepath.Separator)
}
func ListSeparator() string {
	return string(filepath.ListSeparator)
}

func CopyDir(source, dest string) (err error) {

	// get properties of source dir
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// create dest dir

	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}
	directory, _ := os.Open(source)
	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		sourcefilepointer := source + PathSeparator() + obj.Name()

		destinationfilepointer := dest + PathSeparator() + obj.Name()

		if obj.IsDir() {
			// create sub-directories - recursively
			err = CopyDir(sourcefilepointer, destinationfilepointer)
			if err != nil {
				log.Logf(log.ERROR, "issue copying\n%+v", err)
			}
		} else {
			// perform copy
			err = CopyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				log.Logf(log.ERROR, "issue copying\n%+v", err)
			}
		}
	}
	return
}
func CopyFile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
			log.Logf(log.ERROR, "issue chmod\n%+v", err)
		}
	}
	return
}

func Delete(path string) bool {
	err := os.Remove(path)
	if err == nil {
		return true
	}
	return false
}
func DeleteAll(path string) bool {
	err := os.RemoveAll(path)
	if err == nil {
		return true
	}
	return false
}

// MkDir create a directory if it doesn't exist passing in a FileMode such as os.ModePerm
func MkDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Logf(log.ERROR, "issue making dir\n%+v", err)
		}
	}
}

func Available(path string) bool {

	log.Logln(log.DEBUG, "File:", path)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}

func Name(path string) string {

	idx := strings.LastIndex(path, string(os.PathSeparator))
	rns := []rune(path)
	name := string(rns[idx+1 : len(path)])

	return name
}

func PathOnly(path string) string {
	return filepath.Dir(path)
}

func SizeAs(size int64, szType filesize.SizeTypes, si bool) float64 {

	if size > 0 {
		switch szType {
		case filesize.Kilo:
			return filesize.SizeTypes(filesize.Kilo).Convert(int(size), filesize.Bytes, si)
		case filesize.Mega:
			return filesize.SizeTypes(filesize.Mega).Convert(int(size), filesize.Bytes, si)
		case filesize.Giga:
			return filesize.SizeTypes(filesize.Giga).Convert(int(size), filesize.Bytes, si)
		case filesize.Tera:
			return filesize.SizeTypes(filesize.Tera).Convert(int(size), filesize.Bytes, si)
		}
	}
	return 0
}
func FixPath(path string) string {
	if !filepath.IsAbs(path) {
		pth, _ := filepath.Abs(path)
		return pth
	}
	return path
}
func ExpandPath(filepath string) (expandedPath string) {
	cleanedPath := path.Clean(filepath)
	expandedPath = cleanedPath
	if strings.HasPrefix(cleanedPath, "~/") {
		rest := cleanedPath[2:]
		expandedPath = path.Join(HomeFolder(), rest)
	}
	return
}
func HomeFolder() string {

	home, err := os.UserHomeDir()
	if err != nil {
		log.Logf(log.FATAL, "home dir not defined %+v", err)
	}
	return home
}
func FileAvailable(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}
