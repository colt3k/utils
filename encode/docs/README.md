# encode Documentation

These notes split the repository-wide docs down to the `encode` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/encode`
- Run `cd encode && go test ./...` after changing encoding behavior.
- Keep the generated enum in `encodeenum` in sync with supported encodings.
- Preserve backwards-compatible behavior for `B64DecodeStdSanitized`.

## API Surface
- `Encode(data, encodeenum.Encoding)` converts bytes to a string representation.
- `Decode(data, encodeenum.Encoding)` decodes bytes using the selected encoding.
- `B64DecodeStdSanitized(data string)` trims quoting characters before decoding standard base64.
- `encodeenum.Encoding` lists the supported formats.

## Schema & Data Shapes
- No external file schema is used.
- The main typed shape is `encodeenum.Encoding`.

## Dataflow
1. Select the encoding enum that matches the caller contract.
2. Encode outbound bytes or decode inbound text.
3. Use the sanitized helper when inputs may contain quotes or backticks.

## Operator Notes
- Library-only module with no operator-managed runtime.
- If you change the enum list, review downstream modules such as `hash` and `file` for compatibility.
