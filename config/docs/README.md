# config Documentation

These notes split the repository-wide docs down to the `config` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/config`
- Run `cd config && go test ./...` after changing the loader or logger adapter.
- Assign `Config.Logger` before calling `Load`, or use `Nop`/`NewStandardLogger` to avoid nil-interface panics.
- Keep path-resolution order stable unless you are intentionally changing lookup semantics.

## API Surface
- `Config` holds the `go-up` handle and logger.
- `NewConfig()` allocates an empty config wrapper.
- `(*Config).Load(file string)` resolves and loads the env file.
- `Logger`, `LogFunc`, `Standard`, and `Nop` define logging behavior.

## Schema & Data Shapes
- The module loads key/value env files such as `.env`; it does not define a richer typed schema yet.
- `Save` and `Delete` are placeholders and currently do not persist config state.

## Dataflow
1. Build a `Config` and set its logger implementation.
2. Call `Load` with an explicit file name or allow it to default to `.env`.
3. Consume values from the resulting `go-up` handle in the caller.

## Operator Notes
- This is a library module, not a service.
- If configuration stops loading, inspect the executable directory, caller directory, and current working directory in that order.
