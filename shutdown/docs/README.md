# shutdown Documentation

These notes split the repository-wide docs down to the `shutdown` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/shutdown`
- Run `cd shutdown && go test ./...` when adding tests or changing signal handling.
- Keep singleton initialization semantics stable so callers can share one shutdown context.
- Use only signals that map cleanly across the supported environments.

## API Surface
- `SetupNotifyContext(ctx)` returns the shared `Graceful` instance.
- `Graceful.Graceful(cleanup)` waits for `SIGINT` or `SIGTERM`, runs cleanup, and stops notifications.

## Schema & Data Shapes
- No external schema is used.
- The main runtime shape is `Graceful`, which holds the notify context and cancel function.

## Dataflow
1. Initialize the shutdown helper from a parent context.
2. Share the notify context with workers that should stop on shutdown.
3. Run the cleanup callback when the process receives an interrupt or terminate signal.

## Operator Notes
- This module directly affects graceful shutdown behavior in applications that use it.
- If cleanup does not run, inspect context wiring and signal delivery before changing the helper.
