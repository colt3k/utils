# ctxt Documentation

These notes split the repository-wide docs down to the `ctxt` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/ctxt`
- Run `cd ctxt && go test ./...` when adding behavior or tests.
- Keep helpers cancellation-aware and dependency-light.
- Avoid duplicating generic context utilities that belong in the standard library or caller code.

## API Surface
- `SleepContext(ctx, d)` waits for a duration unless the context is canceled first.

## Schema & Data Shapes
- No external schema is used.
- Inputs are a `context.Context` and a `time.Duration`.

## Dataflow
1. Pass a live context into `SleepContext`.
2. Wait for the duration or cancellation signal.
3. Handle `ctx.Err()` when the sleep is interrupted.

## Operator Notes
- Library-only module with no runtime service.
- Most regressions here show up as shutdown or retry timing bugs in downstream modules.
