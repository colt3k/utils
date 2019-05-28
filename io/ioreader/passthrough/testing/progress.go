package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/colt3k/utils/mathut"

	"github.com/colt3k/utils/io/ioreader/passthrough"
)

func main() {
	ctx := context.Background()

	fin, err := os.Open("/Users/username/Downloads/Win10/win10.iso")
	if err != nil {
		panic(err)
	}
	r := passthrough.NewReader(fin)
	fio, err := fin.Stat()
	if err != nil {
		panic(err)
	}
	size := fio.Size()

	//Start a goroutine printing progress
	go func() {
		progressChan := passthrough.NewTicker(ctx, r, int64(size), 1*time.Microsecond)
		for p := range progressChan {
			fmt.Printf("\r%v remaining at %v%% ...", p.Remaining().Round(time.Second), mathut.FmtFloatWithPrecision(p.Percent(), 2))
		}
		fmt.Println("\rdownload is completed")
	}()

	fw, err := os.OpenFile("win10.iso", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	// use the Reader as normal
	if _, err := io.Copy(fw, r); err != nil {
		log.Fatalln(err)
	}
}
