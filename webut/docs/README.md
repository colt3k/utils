# webut Documentation

These notes split the repository-wide docs down to the `webut` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/webut`
- Run `cd webut && go test ./...` if you add tests or expand decoding behavior.
- Keep the package focused on light-weight interoperability helpers rather than full web clients.
- Preserve backwards-compatible boolean parsing for callers that depend on string-encoded booleans.

## API Surface
- `ConvertibleBoolean` accepts JSON booleans encoded either as native booleans or quoted strings.
- `(*ConvertibleBoolean).UnmarshalJSON` implements the tolerant parsing logic.

## Schema & Data Shapes
- The relevant schema is inconsistent JSON boolean fields such as `true`, `false`, `"true"`, or `"false"`.
- The module does not define any external config file format.

## Dataflow
1. Embed `ConvertibleBoolean` in a request or response struct.
2. Decode the JSON payload normally through `encoding/json`.
3. Let the custom unmarshal logic normalize the boolean representation.

## Operator Notes
- Library-only module with no operator-managed runtime.
- If decoding fails, capture the raw upstream payload first to confirm the boolean representation really changed.
