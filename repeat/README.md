# repeat

Fixed-interval task execution with context cancellation support.

## API Surface

- `Rule` configures max attempts plus initial and repeat intervals
- `NewRule()` returns a default rule
- `Process(ctx, rule, processName, repeat, stop)` runs the repeating task until completion or cancellation
- `LoopId()` returns the timestamp-based identifier of the latest loop iteration

## Notes

Use `repeat` when you want scheduled retries with a constant interval. Use `retry` when you want backoff-oriented retry behavior.

## Development

- `cd repeat && go test ./...`
