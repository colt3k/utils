package retry

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Task func() error

type Rule struct {
	MaxAttempts    uint
	currentAttempt uint
	MaxInterval    time.Duration // 60 sec
	MaxElapsed     time.Duration // 15 min
	Elapsed        time.Duration
}

func NewNextBackoff() Rule {
	n := Rule{
		MaxAttempts: 10,
		MaxInterval: 5 * time.Minute,
		MaxElapsed: 15 * time.Minute,
	}
	return n
}
func (r *Rule) NextBackoff() time.Duration {

	jitter := rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(int64(1000))
	pow := math.Pow(float64(2), float64(r.currentAttempt))
	min := math.Min(pow+float64(jitter), float64(r.MaxInterval))
	d := time.Duration(min) * time.Millisecond

	return d
}

func Process(o Task, r Rule) error {
	r.currentAttempt = 1
	if r.MaxInterval == 0 {
		r.MaxInterval = 5 * time.Minute
	}
	if r.MaxAttempts == 0 && r.MaxElapsed == 0 {
		fmt.Println("note: no max attempts or max elapsed time has been set, this will continue until success")
	}
	for {
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
		r.Elapsed+=d
		time.Sleep(d)
		r.currentAttempt++
	}
}