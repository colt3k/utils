# ctxt

Context-aware timing helpers.

## API Surface

`SleepContext(ctx, d)` waits for the requested duration unless the context is canceled first, in which case it returns `ctx.Err()`.

## Notes

Use this package instead of `time.Sleep` when the caller needs cancellation or shutdown awareness.
