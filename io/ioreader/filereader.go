package ioreader

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	ers "github.com/colt3k/nglog/ers/bserr"
	log "github.com/colt3k/nglog/ng"
	"github.com/iancoleman/orderedmap"

	cio "github.com/colt3k/utils/io"
	"github.com/colt3k/utils/io/data"
	"github.com/colt3k/utils/io/ioreader/iocsv"
)

type FileReader struct {
	file string
	cio.Reader
}

func NewFileReader(path string) *FileReader {
	return &FileReader{file: path}
}

/*
LineScanner does not deal well with lines longer than 65536 characters
In that case set the scanner.Buffer([], MAXSIZE) to a larger value

MaxScanTokenSize = 64 * 1024
startBufSize = 4096 // Size of initial allocation for buffer.

@maxBuf 4096
*/
func (r *FileReader) LineScanner(maxBuf, maxScanTokenSize int) (string, error) {

	if f, err := openFile(r.file); err == nil {
		defer f.Close()
		if maxBuf == -1 {
			maxBuf = 4096
		}
		if maxScanTokenSize == -1 {
			maxScanTokenSize = 64 * 1024
		}
		var buffer bytes.Buffer
		scanner := bufio.NewScanner(f)

		buf := make([]byte, maxBuf)
		scanner.Buffer(buf, maxScanTokenSize)
		for scanner.Scan() {
			//fmt.Println(scanner.Text())
			buffer.WriteString(scanner.Text())
			buffer.WriteString("\n")
		}

		if err := scanner.Err(); err != nil {
			log.Logf(log.FATAL, "issue scanning\n%+v", err)
			return "", err
		}
		return buffer.String(), nil
	}
	return "", errors.New("reading failed")

}

// AsString read file as a string
func (r *FileReader) AsString() (string, error) {
	log.Logln(log.DEBUG, "ToOpen:", r.file)
	f, err := openFile(r.file)
	if err != nil {
		log.Logf(log.FATAL, "issue opening file\n%+v", err)
	}
	defer f.Close()

	rdr := bufio.NewReader(f)
	str, err := cio.ReadLine(rdr)
	if ers.NotErr(err, "filereader: read line error") {
		return str, nil
	}
	return "", errors.New("unable to open file")
}
func (r *FileReader) ReadCSVFromFile(filePath string, skipHeader bool, headAr []string) (*data.Table, error) {
	return iocsv.ReadCSVFromFile(r.file, skipHeader, headAr)
}

// AsCSVIntoMap read csv file into map[string]string
func (r *FileReader) AsCSVIntoMap() *map[string]string {
	return iocsv.ReadKV(r.file)
}

// AsCSVIntoOrderedMap read csv into an ordered map
func (r *FileReader) AsCSVIntoOrderedMap() (*orderedmap.OrderedMap, error) {
	return iocsv.ReadOrderedKV(r.file)
}

//AsBytes read file into []byte
func (r *FileReader) AsBytes() ([]byte, error) {
	dat, err := ioutil.ReadFile(r.file)
	if ers.NotErr(err, "filereader: read file error") {
		return dat, nil
	}
	return nil, err
}

//Bytes read file into passed byte array
func (r *FileReader) Bytes(buf []byte) {
	if f, err := os.Open(r.file); err != nil {
		defer f.Close()
		_, err := f.Read(buf)
		if err != nil {
			log.Logf(log.ERROR, "issue reading bytes %+v", err)
		}
	}
}
func openFile(fileStr string) (*os.File, error) {
	fo, err := os.Open(fileStr)
	if err != nil {
		return nil, err
	}

	return fo, nil
}
func (r *FileReader) URL(url string) error {
	// Create the file
	out, err := os.Create(r.file)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil

}
