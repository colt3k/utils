package shutdown

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
)

func Graceful(context context.Context, cleanup func()) context.Context {
	ctx, stop := signal.NotifyContext(context,
		syscall.SIGINT,  // interrupt: stopped by Ctrl + C
		syscall.SIGTERM, // graceful kill
	)
	go func() {
		s := <-ctx.Done()
		fmt.Println("Got signal:", s)
		fmt.Println("\r- Ctrl+C or (kill -9) received")
		cleanup()
		stop()
	}()
	return ctx
}
