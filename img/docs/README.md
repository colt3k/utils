# img Documentation

These notes split the repository-wide docs down to the `img` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/img`
- Run `cd img && go test ./...` if you add tests or new format handling.
- Keep the package narrow and focused on decoding rather than image processing.
- If you add support for more formats, register decoders via standard-library imports.

## API Surface
- `Image.GetAsRGB(reader io.Reader)` decodes the stream and returns an `*image.RGBA` plus the detected format.

## Schema & Data Shapes
- No config or file schema is used.
- The output shape is a standard-library RGBA image and the input format name.

## Dataflow
1. Pass an image reader into `GetAsRGB`.
2. Decode the image using the registered image decoders.
3. Use the returned RGBA buffer in the caller.

## Operator Notes
- Library-only module with no operational runbook.
- Validate on representative image samples if format support changes.
