# shutdown

Graceful process shutdown helpers built on `signal.NotifyContext`.

## API Surface

- `SetupNotifyContext(ctx)` returns a singleton `Graceful` helper
- `Graceful.Graceful(cleanup)` blocks for `SIGINT` or `SIGTERM`, runs cleanup, and then releases signal resources

## Notes

This package centralizes shutdown wiring so callers can share the same signal-aware context across goroutines and cleanup paths.
