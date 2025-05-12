package io

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/mattn/go-isatty"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type Writer interface {
	WriteOut(data []byte, filePath string)
	WriteOutStr(data, filePath string)
	WriteOutString(filePath string) *os.File
	WriteTempFileOfSize(filesize int64, fileprefix string) (fileName string, fileSize int64)
}

// WriteOut write out []byte data to designated file path
func WriteOut(data []byte, filePath string) (int, error) {
	f, errOpen := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if errOpen != nil {
		return 0, fmt.Errorf("ERROR: opening \n%+v", errOpen)
	}
	w := bufio.NewWriter(f)
	n, errWrite := w.Write(data)
	if errWrite != nil {
		return 0, fmt.Errorf("ERROR: write out file\n%+v", errWrite)
	}
	errFlush := w.Flush()
	if errFlush != nil {
		return 0, fmt.Errorf("ERROR: flushing\n%+v", errFlush)
	}
	errClose := f.Close()
	if errClose != nil {
		return 0, fmt.Errorf("ERROR: closing\n%+v", errClose)
	}
	return n, nil
}

// WriteOutAppend write out and append []byte data to designated file path
func WriteOutAppend(data []byte, filePath string) {
	f, errOpenFile := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if errOpenFile != nil {
		log.Printf("ERROR: openfile %v\n", errOpenFile)
	}
	w := bufio.NewWriter(f)
	_, errWrite := w.Write(data)
	if errWrite != nil {
		log.Printf("ERROR: write out file %v\n", errWrite)
	}
	errFlush := w.Flush()
	if errFlush != nil {
		log.Printf("ERROR: flushing file %v\n", errFlush)
	}
	errClose := f.Close()
	if errClose != nil {
		log.Printf("ERROR: closing file %v\n", errClose)
	}
}

// WriteOutStr write out string data to designated file path
func WriteOutStr(data, filePath string) (int, error) {
	return WriteOut([]byte(data), filePath)
}

// WriteOutString send file to create and returns File object to use
func WriteOutString(filePath string) *os.File {
	path, errAbs := filepath.Abs(filePath)
	if errAbs != nil {
		log.Printf("ERROR: determine abs path %v\n", errAbs)
	}
	f, errCreate := os.Create(path)
	if errCreate != nil {
		log.Printf("ERROR: create file %v\n", errCreate)
	}

	return f
}

func WriteTempFileOfSize(filesize int64, fileprefix string) (fileName string, fileSize int64) {
	hash := sha256.New()
	f, _ := os.CreateTemp("", fileprefix)
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
		// return terminal.IsTerminal(int(v.Fd()))
		return isatty.IsTerminal(v.Fd())
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
