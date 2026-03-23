# debug

Runtime debug helpers for dumping pprof data.

## Available Snapshots

- goroutine stack traces
- heap profiles
- thread creation traces
- block profiles

## API Surface

Each snapshot can be written either to standard error or to a caller-provided writer through functions such as `Stack`, `Heap`, `ThreadCreate`, and `Block`.

## Notes

Use these helpers when you need lightweight diagnostics in an application without wiring a full pprof HTTP endpoint.
