# regx Documentation

These notes split the repository-wide docs down to the `regx` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/regx`
- Run `cd regx && go test ./...` after changing patterns or validators.
- Keep pattern names stable because callers use these helpers as semantic shortcuts.
- Add test cases for every new pattern branch or country/format variant.

## API Surface
- `Find` and `Match` expose generic regex application helpers.
- `CC`, `DATE`, `ISBN`, `POSTALCODE`, and `PRICE` expose named validation/extraction helpers.
- `patterns.go` holds the underlying regex catalog.

## Schema & Data Shapes
- No external schema file is defined.
- The effective schema is the set of named regex constants and the string formats they validate.

## Dataflow
1. Choose the helper that matches the data type you need to validate or extract.
2. Pass the raw input string into the helper.
3. Consume the matched value or boolean result in the caller.

## Operator Notes
- Library-only module with no runtime operator steps.
- Regex changes are easy to get wrong; confirm both positive and negative cases before merging.
