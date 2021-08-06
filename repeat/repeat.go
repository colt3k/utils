package repeat

import (
	"context"
	"fmt"
	"time"
)

type Task func() error

type Rule struct {
	MaxAttempts	uint
	currentAttempt uint
	InitialTimer time.Duration
	RepeatTimer time.Duration
}

func NewRule() Rule {
	n := Rule{
		MaxAttempts: 3,
		InitialTimer: 1 * time.Second,
		RepeatTimer: 3 * time.Second,
	}
	return n
}
func Process(ctx context.Context, r Rule, processName string, repeat Task, stop Task) error {
	var count uint = 0
	if r.MaxAttempts == 0 {
		fmt.Printf("%v - note: no max attempts has been set, this will continue forever\n", processName)
	}
	timer := time.NewTimer(r.InitialTimer)
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("repeater cancelled: %v\n", processName)
			if stop != nil {
				if err := stop(); err != nil {
					return err
				}
			}
			return nil
		case t := <-timer.C:
			if r.MaxAttempts > 0 && count >= r.MaxAttempts {
				timer.Stop()
				return nil
			}
			fmt.Printf("timer fired: %v - %v\n", processName, t)
			if err := repeat(); err != nil {
				return err
			}
			count++
			timer.Reset(r.RepeatTimer)
		}
	}
}