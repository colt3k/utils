# stringut Documentation

These notes split the repository-wide docs down to the `stringut` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/stringut`
- Run `cd stringut && go test ./...` after changing regex-driven or validation behavior.
- Group new helpers carefully; this package is broad, but additions should still be clearly reusable.
- Update tests for both positive and negative validation cases.

## API Surface
- `ToChar*`, `ToBool`, and `ToInt` cover conversion helpers.
- `Reverse`, `HRByteCount`, and substring helpers cover transformation helpers.
- `Extract*` and `Validate*` helpers cover regex-driven extraction and validation.
- `Replace*`, `ContainsOnlyNumeric`, and `Eq` cover normalization and comparison helpers.

## Schema & Data Shapes
- No external schema file is used.
- Inputs and outputs are plain strings, booleans, ints, or slices of strings.

## Dataflow
1. Choose the conversion, extraction, validation, or normalization helper that matches the caller need.
2. Pass the raw string into the helper.
3. Consume the scalar, slice, or boolean result in the caller.

## Operator Notes
- Library-only module; operator concerns are indirect through consuming applications.
- Regex and validation changes can break user input handling quickly, so keep regression coverage broad.
