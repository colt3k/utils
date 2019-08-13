package passthrough

import (
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/colt3k/utils/file"
)

// PassThru wraps an existing io.Reader and forwards the Read() call, while displaying
type PassThru struct {
	rc       io.ReadCloser
	name     string
	partId   int
	fullSize int64

	lock  sync.RWMutex // protects total and err
	total int64        // Total # of bytes transferred
	err   error

	ticker *time.Ticker
	readOnce bool
}

// NewPassThru creates an instance of our PassThru object
func New(readCloser io.ReadCloser, f file.File, notifyInSecs int) *PassThru {
	if notifyInSecs == -1 || notifyInSecs == 0 {
		notifyInSecs = 30
	}
	var fsize int64
	var name string
	if f != nil {
		fsize = f.Size()
		name = f.Name()
	}

	ticker := time.NewTicker(time.Duration(notifyInSecs) * time.Second)
	return &PassThru{rc: readCloser, ticker: ticker, fullSize: fsize, name: name}
}

func NewStream(readCloser io.ReadCloser, fsize int64, name string, partId int, notifyInSecs int) *PassThru {
	if notifyInSecs == -1 || notifyInSecs == 0 {
		notifyInSecs = 30
	}

	ticker := time.NewTicker(time.Duration(notifyInSecs) * time.Second)
	return &PassThru{rc: readCloser, ticker: ticker, fullSize: fsize, name: name, partId:partId}
}

// Read 'overrides' the underlying io.Reader's Read method, used to track byte counts and forward the call.
func (pt *PassThru) Read(p []byte) (n int, err error) {
	n, err = pt.rc.Read(p)
	pt.lock.Lock()
	pt.total += int64(n)
	pt.err = err
	pt.lock.Unlock()

	if err == nil {
		if !pt.readOnce {
			fmt.Println("Starting part #", pt.partId)
			go func() {
				for range pt.ticker.C {
					fmt.Printf("Part %d for %s byte(s) sent %d of %d %s%%\n", pt.partId, pt.name, pt.total, pt.fullSize, strconv.FormatFloat(float64(pt.total)*100/float64(pt.fullSize), 'f', 2, 64))
				}
			}()
			pt.readOnce = true
		}

		if pt.total == pt.fullSize {
			pt.ticker.Stop()
			//fmt.Println(pt.name, "File Read Completed FULLSIZE Match ", pt.total)
		}
	} else if err == io.EOF {
		pt.ticker.Stop()
		//fmt.Println(pt.name, "File Read Completed EOF ", pt.total)
	} else {
		fmt.Println(err)
	}

	return n, err
}

// Close
func (pt *PassThru) Close() error {
	if pt.rc != nil {
		pt.rc.Close()
	}
	if pt.ticker != nil {
		pt.ticker.Stop()
		pt.ticker = nil
	}

	pt.rc = nil
pt.total = -1
pt.fullSize = -1

	//we don't actually have to do anything here, since the buffer is just some data in memory
	//and the error is initialized to no-error
	return nil
}

// N gets the total read so far
func (pt *PassThru) N() int64 {
	var n int64
	pt.lock.RLock()
	n = pt.total
	pt.lock.RUnlock()
	return n
}
func (pt *PassThru) Err() error {
	var err error
	pt.lock.RLock()
	err = pt.err
	pt.lock.RUnlock()
	return err
}
// Len returns the number of bytes of the unread portion of the
// slice.
func (pt *PassThru) Len() int {

	if pt.total >= int64(pt.fullSize) {
		return 0
	}
	return int(int64(pt.fullSize) - pt.total)
}

