package io

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	ers "github.com/colt3k/nglog/ers/bserr"
	"golang.org/x/crypto/ssh/terminal"
)

type Writer interface {
	WriteOut(data []byte, filePath string)
	WriteOutStr(data, filePath string)
	WriteOutString(filePath string) *os.File
	WriteTempFileOfSize(filesize int64, fileprefix string) (fileName string, fileSize int64)
}

// WriteOut write out []byte data to designated file path
func WriteOut(data []byte, filePath string) (int,error) {
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return 0,fmt.Errorf("error opening \n%+v", err)
	}
	w := bufio.NewWriter(f)
	n, err := w.Write(data)
	if err != nil {
		return 0,fmt.Errorf("write out file error\n%+v", err)
	}
	err = w.Flush()
	if err != nil {
		return 0,fmt.Errorf("error flushing\n%+v", err)
	}
	err = f.Close()
	if err != nil {
		return 0,fmt.Errorf("error closing\n%+v", err)
	}
	return n,nil
}
// WriteOut write out []byte data to designated file path
func WriteOutAppend(data []byte, filePath string) {
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	ers.NotErr(err, "openfile: error opening")
	w := bufio.NewWriter(f)
	_, err = w.Write(data)
	ers.NotErr(err, "writer: write out file error")
	w.Flush()
	f.Close()
}

// WriteOutStr write out string data to designated file path
func WriteOutStr(data, filePath string) (int,error) {
	return WriteOut([]byte(data), filePath)
}

// WriteOutString send file to create and returns File object to use
func WriteOutString(filePath string) *os.File {
	path, err := filepath.Abs(filePath)
	ers.NotErr(err, "writer: determine abs path error")
	f, err := os.Create(path)
	ers.NotErr(err, "writer: create file error")

	return f
}

func WriteTempFileOfSize(filesize int64, fileprefix string) (fileName string, fileSize int64) {
	hash := sha256.New()
	f, _ := ioutil.TempFile("", fileprefix)
	ra := rand.New(rand.NewSource(time.Now().UnixNano()))
	defer f.Close()
	writer := io.MultiWriter(f, hash)
	written, _ := io.CopyN(writer, ra, filesize)
	fileName = f.Name()
	fileSize = written
	return
}

func CheckIfTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return terminal.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}

func AppendKeyValue(b *bytes.Buffer, key string, value interface{}, quoteEmptyField bool) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	b.WriteString(key)
	b.WriteByte('=')
	AppendValue(b, value, quoteEmptyField)
}

func AppendValue(b *bytes.Buffer, value interface{}, quoteEmptyField bool) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if !NeedsQuoting(stringVal, quoteEmptyField) {
		b.WriteString(stringVal)
	} else {
		b.WriteString(fmt.Sprintf("%q", stringVal))
	}
}
func NeedsQuoting(text string, quoteEmptyField bool) bool {
	if quoteEmptyField && len(text) == 0 {
		return true
	}
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}
