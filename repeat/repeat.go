package repeat

import (
	"context"
	"time"

	log "github.com/colt3k/nglog/ng"
)

var (
	loopId int64 = 0
)

type Task func() error

type Rule struct {
	MaxAttempts    uint
	currentAttempt uint
	InitialTimer   time.Duration
	RepeatTimer    time.Duration
}

func NewRule() Rule {
	n := Rule{
		MaxAttempts:  3,
		InitialTimer: 1 * time.Second,
		RepeatTimer:  3 * time.Second,
	}
	return n
}
func LoopId() int64 {
	return loopId
}
func Process(ctx context.Context, r Rule, processName string, repeat Task, stop Task) error {
	curTime := time.Now().Format(time.RFC3339)
	log.Logf(log.DEBUG, "setting up repeat process as of %v", curTime)
	var count uint = 0
	if r.MaxAttempts == 0 {
		log.Logf(log.INFO, "%v - note: no max attempts has been set, this will continue forever", processName)
	}
	timer := time.NewTimer(r.InitialTimer)
	for {
		select {
		case <-ctx.Done():
			timer.Stop()
			log.Logf(log.INFO, "repeater cancelled: %v", processName)

			if stop != nil {
				if err := stop(); err != nil {
					return err
				}
			}
			return nil
		case t := <-timer.C:
			// Stop the timer from being run again until our function finishes and returns
			timer.Stop()
			if r.MaxAttempts > 0 && count >= r.MaxAttempts {
				return nil
			}

			loopId = time.Now().Unix()
			log.Logf(log.INFO, "%v timer fired: %v - %v", loopId, processName, t)
			if err := repeat(); err != nil {
				log.Logf(log.ERROR, "(in repeater) %v exited task iteration with error %v", loopId, err)
				return err
			}
			log.Logf(log.INFO, "(in repeater) %v exited task iteration without error", loopId)
			count++
			timer = time.NewTimer(r.RepeatTimer)
		}
	}
}
