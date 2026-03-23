# debug Documentation

These notes split the repository-wide docs down to the `debug` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/debug`
- Run `cd debug && go test ./...` if you add tests or change writer behavior.
- Keep the package focused on pprof dump helpers rather than higher-level observability wiring.
- Default writer behavior should continue to target stderr when the caller passes nil.

## API Surface
- `Stack`, `Heap`, `ThreadCreate`, and `Block` write named pprof snapshots to a writer.
- `PrintStack`, `PrintHeap`, `PrintThreadCreate`, and `PrintBlock` are stderr convenience wrappers.

## Schema & Data Shapes
- No external schema is used.
- The only variable input is the optional `io.Writer` passed into the dump helpers.

## Dataflow
1. Pick the pprof view you need for the failure mode.
2. Call the matching helper with a writer or allow stderr to be used.
3. Inspect the resulting dump in logs or capture it for offline analysis.

## Operator Notes
- Library-only module; operator use is usually ad hoc during incident debugging.
- Do not call the dump helpers in tight loops unless the caller explicitly wants large diagnostic output.
