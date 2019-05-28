package concur

import (
	"sync"
	"sync/atomic"
)

type ResyncOnce struct {
	m    sync.Mutex
	done uint32
}

func (o *ResyncOnce) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 1 {
		return
	}
	// Slow-path.
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}

func (o *ResyncOnce) Reset() {
	if atomic.LoadUint32(&o.done) == 1 {
		o.m.Lock()
		defer o.m.Unlock()
		if o.done == 1 {
			defer atomic.StoreUint32(&o.done, 0)
		}
	}
}
