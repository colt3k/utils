# qrand Documentation

These notes split the repository-wide docs down to the `qrand` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/qrand`
- Run `cd qrand && go test ./...` when changing network or decoding behavior.
- Keep the external API contract visible because this module depends on a remote public service.
- Request at least 8 bytes when generating a `uint64` seed.

## API Surface
- `GenerateSeedData(amount int)` fetches random bytes and converts them into a big-endian `uint64` seed.
- `Response` matches the QRNG JSON payload.

## Schema & Data Shapes
- The remote API returns JSON with `type`, `length`, `data`, and `success` fields.
- The module itself does not define a local config file schema.

## Dataflow
1. Build the QRNG request URL with the desired byte length.
2. Fetch and decode the JSON response.
3. Convert the returned byte slice into the seed value for caller use.

## Operator Notes
- This module depends on external network availability.
- If seed generation fails, verify outbound access to `https://qrng.anu.edu.au` before changing code.
