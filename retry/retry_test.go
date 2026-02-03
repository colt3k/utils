package retry

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var (
	counter = 0
)

func secondRule() Rule {
	return Rule{
		MaxAttempts:             3,
		MaxInterval:             10,
		MaxIntervalDurationType: time.Second,
		MaxElapsed:              15 * time.Second,
		SleepDurationType:       Seconds,
	}
}
func minuteRule() Rule {
	return Rule{
		MaxAttempts:             3,
		MaxInterval:             1,
		MaxIntervalDurationType: time.Minute,
		MaxElapsed:              10 * time.Minute,
		SleepDurationType:       Minutes,
	}
}
func millisecondRule() Rule {
	return Rule{
		MaxAttempts:             3,
		MaxInterval:             1000,
		MaxIntervalDurationType: time.Millisecond,
		MaxElapsed:              15000 * time.Millisecond,
		SleepDurationType:       Milliseconds,
	}
}
func TestRetry(t *testing.T) {
	err := Process(context.Background(), func() error {
		return doSomething("somevalue")
	}, NewRule())
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

func TestMillisecondRetry(t *testing.T) {
	err := Process(context.Background(), func() error {
		return doSomething("somevalue")
	}, millisecondRule())
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

func TestSecondRetry(t *testing.T) {
	err := Process(context.Background(), func() error {
		return doSomething("somevalue")
	}, secondRule())
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

// run two times first time fails, second succeeds to stop retry
func doSomething(someparam string) error {
	counter++
	jitter := rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(int64(5)) //nolint:gosec
	fmt.Printf("someparam %v, called %v, rand %v\n", someparam, counter, jitter)

	if jitter == 2 {
		return nil
	}
	return fmt.Errorf("issue %v", counter)
}
