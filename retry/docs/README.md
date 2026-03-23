# retry Documentation

These notes split the repository-wide docs down to the `retry` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/retry`
- Run `cd retry && go test ./...` after changing retry or backoff behavior.
- Keep defaulting rules stable because callers often rely on omitted fields.
- Cover cancellation, elapsed-time ceilings, and duration-type interactions in tests.

## API Surface
- `Rule` defines max attempts, interval settings, elapsed ceilings, and duration units.
- `NewRule()` returns the default retry rule.
- `Process(ctx, task, rule)` executes the retry loop.
- `NextBackoff` calculates the next backoff duration.

## Schema & Data Shapes
- No external file schema is defined.
- The important runtime shape is `Rule`, especially `MaxAttempts`, `MaxInterval`, `MaxIntervalDurationType`, `MaxElapsed`, and `SleepDurationType`.

## Dataflow
1. Build a retry `Rule` with the desired limits and units.
2. Call `Process` with the target task and context.
3. Repeat the task until it succeeds, the rule is exhausted, or the context is canceled.

## Operator Notes
- Library-only module; issues surface as stalled retries or too-aggressive retry storms in consumers.
- When behavior looks wrong, inspect unit mismatches between `MaxInterval` and `SleepDurationType` first.
