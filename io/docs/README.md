# io Documentation

These notes split the repository-wide docs down to the `io` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/io`
- Run `cd io && go test ./...` after changing read/write behavior.
- Keep subpackage responsibilities clear: `iocsv`, `ioexif`, `ioimage`, `iotab`, `passthrough`, and `logstash` each serve a distinct job.
- Be careful with stream semantics because other modules depend on these helpers in the middle of data pipelines.

## API Surface
- `Open`, `ReadLine`, `AsBytes`, `AsString`, and `LastLineWithSeek` provide general reader helpers.
- `WriteOut`, `WriteOutAppend`, `WriteTempFileOfSize`, and quoting helpers provide general writer behavior.
- `iocsv.ReadCSV*`, `ioexif.New`, `ioimage.NewImageMeta`, and `iotab.ReadTabDelim` cover format-specific parsing.
- `passthrough.New` and related types wrap streams with progress reporting.

## Schema & Data Shapes
- No single external schema is enforced by this module.
- Main data shapes include `data.Table`, ordered key/value maps from `iocsv`, and pass-through progress records.

## Dataflow
1. Open or create the source stream.
2. Transform it through the format-specific reader or pass-through wrapper.
3. Persist or forward the resulting bytes, rows, or metadata to the caller.

## Operator Notes
- Library-only module; issues here usually show up as truncated reads, malformed parsing, or missing progress output downstream.
- Re-test the consuming module when changing CSV, EXIF, image, or pass-through behavior.
