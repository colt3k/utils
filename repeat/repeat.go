package repeat

import (
	"context"
	"fmt"
	"time"
)

type Task func() error

func Process(ctx context.Context, initTimer, repeatTimer int, o Task, stop Task) error {
	timer := time.NewTimer(time.Second * time.Duration(initTimer))
	for {
		select {
		case <-ctx.Done():
			if stop != nil {
				fmt.Println("Cleanup started")
				if err := stop(); err != nil {
					return err
				}
				fmt.Println("Cleanup complete")
			}
			return nil
		case t := <-timer.C:
			fmt.Printf("timer fired: %v\n", t)
			if err := o(); err != nil {
				return err
			}
			timer.Reset(time.Second * time.Duration(repeatTimer))
		}
	}
}