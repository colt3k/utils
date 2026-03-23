# repeat Documentation

These notes split the repository-wide docs down to the `repeat` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/repeat`
- Run `cd repeat && go test -mod=mod ./...` in this repository state because vendoring is out of sync.
- Keep timer lifecycle explicit to avoid goroutine and timer leaks.
- Use `retry` for backoff-oriented behavior; use this module for constant-interval loops.

## API Surface
- `Rule` configures max attempts plus initial and repeat timers.
- `NewRule()` returns the default rule.
- `LoopId()` exposes the timestamp-based id of the latest loop iteration.
- `Process(ctx, rule, processName, repeat, stop)` runs the repeating workflow.

## Schema & Data Shapes
- No external schema file is used.
- The main runtime shapes are `Rule` and the repeat/stop callback functions.

## Dataflow
1. Create or customize a `Rule`.
2. Call `Process` with the repeat function and optional stop hook.
3. Let the loop continue until max attempts, cancellation, or a callback error ends it.

## Operator Notes
- Library-only module; operational issues usually show up as stuck or overly chatty loops in consumers.
- When loops misbehave, inspect timer values, cancellation paths, and callback error handling first.
