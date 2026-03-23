// Package shutdown centralizes SIGINT and SIGTERM handling.
package shutdown

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"
)

var once sync.Once
var (
	instance *Graceful
)

type Graceful struct {
	NotifyContext context.Context
	stop          context.CancelFunc
}

// SetupNotifyContext initializes or returns the shared signal-aware context.
func SetupNotifyContext(ctxt context.Context) *Graceful {
	once.Do(func() {
		instance = &Graceful{}
		ctx, stop := signal.NotifyContext(ctxt,
			syscall.SIGINT,  // interrupt: stopped by Ctrl + C
			syscall.SIGTERM, // graceful kill
		)
		instance.NotifyContext = ctx
		instance.stop = stop
	})

	return instance
}

// Graceful waits for shutdown, runs cleanup, and then releases resources.
func (g Graceful) Graceful(cleanup func()) {
	s := <-g.NotifyContext.Done()
	fmt.Println("Got signal:", s)
	fmt.Println("\r- Ctrl+C or (kill -9) received")
	cleanup()
	g.stop()
}
