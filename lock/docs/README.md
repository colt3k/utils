# lock Documentation

These notes split the repository-wide docs down to the `lock` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/lock`
- Run `cd lock && go test ./...` after changing lock acquisition semantics.
- Keep the implementation lightweight and based on exclusive file creation.
- Preserve the temp-directory naming convention unless downstream callers are updated too.

## API Surface
- `New(appName)` creates a lock path in the system temp directory.
- `Try()` acquires the lock by creating the file exclusively.
- `Unlock()` removes the file when work is complete.

## Schema & Data Shapes
- No external schema is used.
- The lock state is represented by a single lock file named `<appName>.lck` in the temp directory.

## Dataflow
1. Create a lock instance for the application name.
2. Call `Try` before doing work that must be single-instance.
3. Call `Unlock` during orderly shutdown or cleanup.

## Operator Notes
- Library-only module; operators mainly care about stale lock files after crashes.
- If a process exits unexpectedly, inspect and remove the temp lock file before retrying.
