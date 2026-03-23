# profile Documentation

These notes split the repository-wide docs down to the `profile` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/profile`
- Run `cd profile && go test ./...` when adding tests or changing metrics output.
- Keep helpers lightweight so callers can drop them into applications without major setup.
- Use `profile/test/main.go` as the example entry point.

## API Surface
- `Duration(invocation, name)` reports elapsed time.
- `MemUsage()` reports heap and allocation statistics.
- `CpuUsagePercent()` reports process CPU usage.

## Schema & Data Shapes
- No external schema is used.
- Outputs are human-readable metrics rather than typed records.

## Dataflow
1. Capture the start time or target process state.
2. Call the relevant profiling helper at the point of interest.
3. Inspect the emitted metrics in logs or console output.

## Operator Notes
- Library-only module; operator usage is normally ad hoc during tuning or debugging.
- Do not treat these helpers as a replacement for full observability when deeper tracing is required.
