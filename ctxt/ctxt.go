// Package ctxt contains small context-aware helpers.
package ctxt

import (
	"context"
	"time"
)

// SleepContext waits for d unless ctx is canceled first.
func SleepContext(ctx context.Context, d time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(d):
		return nil
	}
}
