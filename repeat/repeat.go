package repeat

import (
	"context"
	"fmt"
	"time"
)

type Task func() error

func Process(ctx context.Context, initTimer, repeatTimer int, o Task) error {
	timer := time.NewTimer(time.Second * time.Duration(initTimer))
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Done or Signal")
			return nil
		case t := <-timer.C:
			fmt.Printf("timer fired: %v\n", t)
			if err := o(); err == nil {
				return nil
			}
			timer.Reset(time.Second * time.Duration(repeatTimer))
		}
	}
}