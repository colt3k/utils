package retry

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"
)

type DurationType string

const (
	Milliseconds DurationType = "m"
	Seconds      DurationType = "s"
)

type Task func() error

type Rule struct {
	MaxAttempts    uint
	currentAttempt uint
	MaxInterval    time.Duration // 60 sec
	MaxElapsed     time.Duration // 15 min
	Elapsed        time.Duration
	DurationType   DurationType
}

func NewRule() Rule {
	n := Rule{
		MaxAttempts:  10,
		MaxInterval:  5 * time.Minute,
		MaxElapsed:   15 * time.Minute,
		DurationType: Milliseconds,
	}
	return n
}

func (r *Rule) NextBackoff() time.Duration {
	var jitter int64
	var pow, min float64
	var d time.Duration
	switch r.DurationType {
	case "m":
		jitter = rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(int64(1000)) //nolint:gosec
		pow = math.Pow(float64(2), float64(r.currentAttempt))
		min = math.Min(pow+float64(jitter), float64(r.MaxInterval))
		d = time.Duration(min) * time.Millisecond
	case "s":
		jitter = rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(int64(60)) //nolint:gosec
		pow = math.Pow(float64(2), float64(r.currentAttempt))
		min = math.Min(pow+float64(jitter), float64(r.MaxInterval))
		d = time.Duration(min) * time.Second
	}

	return d
}

func Process(ctx context.Context, o Task, r Rule) error {
	r.currentAttempt = 1
	if r.MaxInterval == 0 {
		r.MaxInterval = 5 * time.Minute
	}
	if r.MaxAttempts == 0 && r.MaxElapsed == 0 {
		fmt.Println("note: no max attempts or max elapsed time has been set, this will continue until success")
	}
	if len(r.DurationType) == 0 {
		r.DurationType = Milliseconds
	}
	// first time execute after 1 millisecond
	timer := time.NewTimer(time.Millisecond * 1)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-timer.C:
			if r.currentAttempt == r.MaxAttempts+1 {
				return fmt.Errorf("exceeded attempts")
			}
			// Only exit if set to something other than 0
			if r.MaxElapsed > 0 && (r.Elapsed > r.MaxElapsed) {
				return fmt.Errorf("exceeded maximum elapsed")
			}
			// if no error then exit
			if err := o(); err == nil {
				return nil
			}
			d := r.NextBackoff()
			fmt.Printf("failed retrying in %v...\n", d)
			r.Elapsed += d
			r.currentAttempt++
			timer.Reset(d)
		}
	}
}
