package file

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"log"

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
				log.Printf("ERROR: issue copying\n%+v\n", err)
			}
		} else {
			// perform copy
			err = CopyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				log.Printf("ERROR: issue copying\n%+v\n", err)
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
			log.Printf("ERROR: issue chmod\n%+v\n", err)
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
			log.Printf("ERROR: issue making dir\n%+v\n", err)
		}
	}
}

func Available(path string) bool {

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
		log.Fatalf("FATAL: home dir not defined %+v", err)
	}
	return home
}
func FileAvailable(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}

type ByNumericalFilename []os.FileInfo

func (nf ByNumericalFilename) Len() int      { return len(nf) }
func (nf ByNumericalFilename) Swap(i, j int) { nf[i], nf[j] = nf[j], nf[i] }
func (nf ByNumericalFilename) Less(i, j int) bool {

	// Use path names
	pathA := nf[i].Name()
	pathB := nf[j].Name()

	// Grab integer value of each filename by parsing the string and slicing off
	// the extension
	re := regexp.MustCompile("[0-9]+")
	var a int64
	var b int64
	var err1 error
	var err2 error
	oneAr := re.FindAllString(pathA, 1)
	twoAr := re.FindAllString(pathB, 1)
	if oneAr != nil {
		a, err1 = strconv.ParseInt(oneAr[0], 10, 64)
	} else {
		err1 = fmt.Errorf("no numbers found")
	}
	if twoAr != nil {
		b, err2 = strconv.ParseInt(twoAr[0], 10, 64)
	} else {
		err2 = fmt.Errorf("no numbers found")
	}

	// If any were not numbers sort lexicographically
	if err1 != nil || err2 != nil {
		return pathA < pathB
	}

	// Which integer is smaller?
	return a < b
}

type ByNumericalFilenameRev []os.FileInfo

func (nf ByNumericalFilenameRev) Len() int      { return len(nf) }
func (nf ByNumericalFilenameRev) Swap(i, j int) { nf[i], nf[j] = nf[j], nf[i] }
func (nf ByNumericalFilenameRev) Less(i, j int) bool {

	// Use path names
	pathA := nf[i].Name()
	pathB := nf[j].Name()

	// Grab integer value of each filename by parsing the string and slicing off
	// the extension
	re := regexp.MustCompile("[0-9]+")
	var a int64
	var b int64
	var err1 error
	var err2 error
	oneAr := re.FindAllString(pathA, 1)
	twoAr := re.FindAllString(pathB, 1)
	if oneAr != nil {
		a, err1 = strconv.ParseInt(oneAr[0], 10, 64)
	} else {
		err1 = fmt.Errorf("no numbers found")
	}
	if twoAr != nil {
		b, err2 = strconv.ParseInt(twoAr[0], 10, 64)
	} else {
		err2 = fmt.Errorf("no numbers found")
	}

	// If any were not numbers sort lexicographically
	if err1 != nil || err2 != nil {
		return pathB < pathA
	}

	// Which integer is smaller?
	return b < a
}

func RotateFiles(fileName string) error {
	parts := strings.Split(filepath.Base(fileName), ".")
	ext := parts[len(parts)-1]
	files, err := ioutil.ReadDir(path.Dir(fileName))
	if err != nil {
		return err
	}

	sort.Sort(ByNumericalFilename(files))
	found := make([]os.FileInfo, 0)
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), "."+ext) {
			fmt.Println(file.Name())
			found = append(found, file)
		}
	}

	sort.Sort(ByNumericalFilenameRev(found))
	// Rename to 1 higher
	re := regexp.MustCompile("[0-9]+")
	for _, j := range found {
		if ar := re.FindAllString(j.Name(), 1); ar != nil {
			lastNum, err := strconv.ParseInt(ar[0], 10, 64)
			if err != nil {
				fmt.Printf("issue renaming convert found int %v\n", err)
			}
			newNum := lastNum + 1
			oldName := j.Name()
			newName := strings.ReplaceAll(oldName, strconv.Itoa(int(lastNum)), strconv.Itoa(int(newNum)))
			err = os.Rename(path.Join(path.Dir(fileName), oldName), path.Join(path.Dir(fileName), newName))
			if err != nil {
				fmt.Printf("issue renaming %v\n", err)
			}
		}

	}

	return nil
}

func findFile(fileName string) ([]os.FileInfo, error) {
	parts := strings.Split(filepath.Base(fileName), ".")
	ext := parts[len(parts)-1]
	re := regexp.MustCompile("[a-zA-Z]+")
	pfx := ""
	if ar := re.FindAllString(filepath.Base(fileName), 1); ar != nil {
		pfx = ar[0]
	}
	files, err := ioutil.ReadDir(path.Dir(fileName))
	if err != nil {
		return nil, err
	}
	sort.Sort(ByNumericalFilename(files))
	found := make([]os.FileInfo, 0)

	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), pfx) && strings.HasSuffix(file.Name(), "."+ext) {
			found = append(found, file)
		}
	}
	return found, nil
}
func CleanRotatedByCount(fileName string, max int) error {
	found, err := findFile(fileName)
	if err != nil {
		return err
	}
	sort.Sort(ByNumericalFilename(found))
	for i, j := range found {
		if i >= max {
			err := os.Remove(path.Join(path.Dir(fileName), j.Name()))
			if err != nil {
				fmt.Printf("issue removing %v\n", err)
			}
			//fmt.Printf("removing %v\n", j.Name())
		}
	}
	return nil
}

func CleanRotatedByDays(fileName string, days int) error {
	t := time.Now()
	found, err := findFile(fileName)
	if err != nil {
		return err
	}
	sort.Sort(ByNumericalFilename(found))
	for _, j := range found {
		hrs := int(t.Sub(j.ModTime()).Hours())
		if (hrs / 24) > days {
			err := os.Remove(path.Join(path.Dir(fileName), j.Name()))
			if err != nil {
				fmt.Printf("issue removing %v\n", err)
			}
			//fmt.Printf("removing %v\n", j.Name())
		}
	}
	return nil
}
