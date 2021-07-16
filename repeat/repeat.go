package repeat

import (
	"context"
	"fmt"
	"time"
)

type Task func() error

func Process(ctx context.Context, initTimer, repeatTimer int, o Task) {
	timer := time.NewTimer(time.Second * time.Duration(initTimer))
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Done or Signal")
			return
		case t := <-timer.C:
			fmt.Printf("timer fired: %v\n", t)
			o()
			timer.Reset(time.Second * time.Duration(repeatTimer))
		}
	}
}