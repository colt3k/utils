package io

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strings"

	log "github.com/colt3k/nglog/ng"
	"github.com/gorilla/http"
	"github.com/iancoleman/orderedmap"

	"github.com/colt3k/utils/io/data"
)

type Reader interface {
	Bytes(buf []byte)
	AsBytes() ([]byte, error)
	AsString() (string, error)
	LineScanner(maxBuf, maxScanTokenSize int) (string, error)
	AsCSVIntoMap() *map[string]string
	AsCSVIntoOrderedMap() (*orderedmap.OrderedMap, error)
	ReadCSVFromFile(filePath string, skipHeader bool, headAr []string) (*data.Table, error)
}

// Open returns an io.ReadCloser representing the contents of the
// source specified by a uri.
func Open(uri string) (io.ReadCloser, error) {
	switch {
	case uri == "-":
		return newStdinReader()
	case strings.HasPrefix(uri, "file://"):
		return newFileReader(uri)
	case strings.HasPrefix(uri, "http://"), strings.HasPrefix(uri, "https://"):
		return newHttpReader(uri)
	case strings.HasPrefix(uri, "tcp://"):
		return newTcpReader(uri)
	}
	return nil, fmt.Errorf("no handler registered for %q", uri)
}

type readCloser struct {
	io.Reader
}

func (r *readCloser) Close() error { return nil }

func newStdinReader() (io.ReadCloser, error) {
	return &readCloser{os.Stdin}, nil
}

func newFileReader(uri string) (io.ReadCloser, error) {
	fname := strings.TrimPrefix(uri, "file://")
	return os.Open(fname)
}

func newHttpReader(uri string) (io.ReadCloser, error) {
	status, _, body, err := http.DefaultClient.Get(uri, nil)
	if err != nil {
		return nil, err
	}
	if !status.IsSuccess() {
		return nil, &http.StatusError{status}
	}
	return body, nil
}

func newTcpReader(uri string) (io.ReadCloser, error) {
	dst := strings.TrimPrefix(uri, "tcp://")
	conn, err := net.Dial("tcp", dst)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

/*
Line returns a single line (without the ending \n)
from the input buffered reader.
An error is returned if` there is an error with the
buffered reader.
*/
func ReadLine(br *bufio.Reader) (string, error) {
	var (
		isPrefix = true
		err      error
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = br.ReadLine()

		ln = append(ln, line...)
		//fmt.Println(string(ln))
	}
	return string(ln), err
}

// Bytes read file into passed byte array
func Bytes(buf []byte) {
}

// AsBytes read file into []byte
func AsBytes() ([]byte, error) {
	return nil, nil
}

// AsString read file as a string
func AsString() (string, error) {
	return "", nil
}

func LineScanner(maxBuf, maxScanTokenSize int) (string, error) {
	return "", nil
}

// AsCSVIntoMap read csv file into map[string]string
func AsCSVIntoMap() *map[string]string {
	return nil
}

func LastLineWithSeek(filepath string, amt int) ([]string, error) {
	fileHandle, err := os.Open(filepath)

	if err != nil {
		return nil, fmt.Errorf("cannot open file %s\n%+v", filepath, err)
	}
	defer fileHandle.Close()

	line := ""
	var cursor int64 = 0
	stat, _ := fileHandle.Stat()
	filesize := stat.Size()
	if filesize <= 0 {
		return nil, fmt.Errorf("empty file %s", filepath)
	} else if filesize <= 2000 {
		// Check if only new lines???
		b, err := ioutil.ReadAll(fileHandle)
		if err != nil {
			return nil, fmt.Errorf("error reading file %s\n%+v", filepath, err)
		}
		re := regexp.MustCompile(`\r?\n`)
		tmpInput := re.ReplaceAllString(string(b), "")
		if len(tmpInput) == 0 {
			err = fileHandle.Close()
			if err != nil {
				log.Logf(log.FATAL, "issue closing file\n%+v", err)
			}
			err = os.Truncate(filepath, 0)
			if err != nil {
				return nil, fmt.Errorf("empty file except newlines error truncating %s\n%+v", filepath, err)
			}
			return nil, fmt.Errorf("empty file except newlines truncated %s", filepath)
		}

	}

	lines := make([]string, 0)
	count := 0

	for {

		cursor -= 1
		_, err = fileHandle.Seek(cursor, io.SeekEnd)
		if err != nil {
			log.Logf(log.ERROR, "issue seeking %+v", err)
		}

		char := make([]byte, 1)
		_, err := fileHandle.Read(char)
		if err != nil {
			log.Logf(log.ERROR, "issue reading %+v", err)
		}

		// 10 line feed, 13 carriage return
		if cursor != -1 && (char[0] == 10 || char[0] == 13) { // stop if we find a line
			rns := []rune(line)
			if rns[len(rns)-1] == 10 || rns[len(rns)-1] == 13 {
				lines = append(lines, string(rns[:len(rns)-1]))
			} else {
				lines = append(lines, line)
			}
			count++
			line = ""
			if count >= amt {
				break
			}
		}

		line = fmt.Sprintf("%s%s", string(char), line) // there is more efficient way

		if cursor == -filesize { // stop if we are at the beginning, send back the only line
			lines = append(lines, line)
			break
		}
	}

	return lines, nil
}
