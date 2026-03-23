# lock

Simple lock mechanism.

## Behavior

`lock.New(appName)` creates a lock file path under the system temp directory using `<appName>.lck`. `Try` acquires the lock with exclusive create semantics, and `Unlock` removes the file when work is complete.

## Notes

This is a lightweight process-coordination helper, not a distributed lock. Callers should clear stale lock files when a previous process terminates unexpectedly.
