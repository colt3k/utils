package retry

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var (
	counter = 0
)
func TestRetry(t *testing.T) {
	err := Retry(func() error {
		return doSomething("somevalue")
	}, NewNextBackoff())
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func doSomething (someparam string) error {
	counter++
	jitter := rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(int64(3))
	fmt.Printf("someparam %v, called %v, rand %v\n", someparam, counter, jitter)

	if jitter == 2 {
		return nil
	}
	return fmt.Errorf("issue %v", counter)
}