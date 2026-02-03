package repeat

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestRepeat(t *testing.T) {
	r := Rule{
		MaxAttempts:  2,
		InitialTimer: 1 * time.Second,
		RepeatTimer:  2 * time.Second,
	}
	err := Process(context.Background(), r, "testme", repeatingTask, cleanTask)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	err = Process(context.Background(), r, "testme", failingRepeatingTask, cleanTask)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

func repeatingTask() error {

	fmt.Println("in my task")
	return nil
}

func failingRepeatingTask() error {

	fmt.Println("in my task")
	return fmt.Errorf("issue in my task")
}

func cleanTask() error {

	fmt.Println("stopping task or cleanup")
	return nil
}
