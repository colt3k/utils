package iowriter

import (
	"os"

	"github.com/colt3k/utils/file"
	"log"

	"github.com/colt3k/utils/io"
)

type FileWriter struct {
	file file.File
}

func NewFileWriter() *FileWriter {
	return &FileWriter{}
}

//WriteOut write []byte out to file path
func (w *FileWriter) WriteOut(data []byte, filePath string) {
	_,err := io.WriteOut(data, filePath)
	if err != nil {
		log.Printf("ERROR: issue writing out %+v\n",err)
	}
}

//WriteOutStr write string out to filepath
func (w *FileWriter) WriteOutStr(data, filePath string) {
	_,err := io.WriteOutStr(data, filePath)
	if err != nil {
		log.Printf("ERROR: issue writing out %+v",err)
	}
}

//WriteOutString create file and return pointer
func (w *FileWriter) WriteOutString(filePath string) *os.File {
	return io.WriteOutString(filePath)
}

func (w *FileWriter) WriteTempFileOfSize(filesize int64, fileprefix string) (fileName string, fileSize int64) {
	return io.WriteTempFileOfSize(filesize, fileprefix)
}
