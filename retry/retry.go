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
	Minutes      DurationType = "mi"
)

type Task func() error

/*
Rule

	MaxAttempts how many times to try before quiting
	MaxInterval maximum amount of time to wait between tries, default 5
	MaxIntervalDurationType time Duration of the MaxInterval, default time.Minute
	MaxElapsed maximum total time before quiting all retries
	Elapsed used for tracking total elapsed time, available to read externally
	SleepDurationType time duration type to sleep between retries, default retry.Milliseconds
*/
type Rule struct {
	MaxAttempts             uint
	currentAttempt          uint
	MaxInterval             int64 // 60 sec
	MaxIntervalDurationType time.Duration
	MaxElapsed              time.Duration // 15 min
	elapsed                 time.Duration
	SleepDurationType       DurationType
}

func NewRule() Rule {
	n := Rule{
		MaxAttempts:             10,
		MaxInterval:             5,
		MaxElapsed:              15 * time.Minute,
		MaxIntervalDurationType: time.Minute,
		SleepDurationType:       Milliseconds,
	}
	return n
}

// NextBackoff determines the next backoff duration using the maximum interval and time duration type and case to the Duration Type specified
func (r *Rule) NextBackoff() time.Duration {
	var jitter int64
	var pow, minimum float64
	var d time.Duration
	jitter = rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(int64(r.MaxInterval)) //nolint:gosec
	pow = math.Pow(float64(2), float64(r.currentAttempt))
	minimum = math.Min(pow+float64(jitter), float64(r.MaxInterval))
	// fmt.Printf("jitter: %v, pow: %v, maxinterval: %v, elapsed: %v, maxelapsed: %v\n", jitter, pow, r.MaxInterval, r.elapsed, r.MaxElapsed)
	switch r.SleepDurationType {
	case "m":
		d = time.Duration(minimum) * time.Millisecond
	case "s":
		d = time.Duration(minimum) * time.Second
	case "mi":
		d = time.Duration(minimum) * time.Minute
	}

	return d
}

// Process starts the retry process passing the context, Task and Rule
func Process(ctx context.Context, o Task, r Rule) error {
	r.currentAttempt = 1
	if r.MaxInterval != 0 && r.MaxIntervalDurationType == 0 {
		return fmt.Errorf("MaxIntervalDurationType not set, set on intended time unit for MaxInterval (time.Millisecond, time.Second, time.Minute)")
	}
	if r.MaxInterval == 0 {
		r.MaxInterval = 5
		r.MaxIntervalDurationType = time.Minute
	}
	if r.MaxAttempts == 0 && r.MaxElapsed == 0 {
		fmt.Println("note: no max attempts or max elapsed time has been set, this will continue until success")
	}

	if len(r.SleepDurationType) == 0 {
		r.SleepDurationType = Milliseconds
	}
	// first time execute after 1 millisecond
	timer := time.NewTimer(time.Millisecond * 1)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-timer.C:
			// if no error then exit
			if err := o(); err == nil {
				return nil
			}
			if r.currentAttempt == r.MaxAttempts+1 {
				return fmt.Errorf("exceeded attempts")
			}
			// Only exit if set to something other than 0
			if r.MaxElapsed > 0 && (r.elapsed >= r.MaxElapsed) {
				return fmt.Errorf("exceeded maximum elapsed")
			}
			d := r.NextBackoff()
			fmt.Printf("failed retrying in %v...\n", d)
			r.elapsed += d
			r.currentAttempt++
			timer.Reset(d)
		}
	}
}
