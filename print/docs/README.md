# print Documentation

These notes split the repository-wide docs down to the `print` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/print`
- Run `cd print && go test ./...` if you add tests or change rendering behavior.
- Keep the wrapper intentionally small; richer table configuration belongs in callers or the underlying library.
- Preserve stdout rendering because downstream CLIs rely on that behavior.

## API Surface
- `Printer` stores headers and rows.
- `(*Printer).TablePrint(border bool)` renders the configured table to stdout.

## Schema & Data Shapes
- No external schema is used.
- The runtime shape is just a header slice plus row slices of strings.

## Dataflow
1. Populate a `Printer` with header and row data.
2. Choose whether borders should be rendered.
3. Call `TablePrint` to emit the formatted table.

## Operator Notes
- Library-only module with no runtime service.
- Validate output formatting in the consuming CLI if table layout changes.
