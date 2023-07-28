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

func TestRetry(t *testing.T) {
	err := Process(context.Background(), func() error {
		return doSomething("somevalue")
	}, NewRule())
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
