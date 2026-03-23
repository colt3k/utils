# qrand

Quantum random seed generation backed by the Australian National University QRNG API.

## API Surface

- `GenerateSeedData(amount int)` requests random bytes and returns them as a big-endian `uint64` seed

## Notes

The package uses the public JSON endpoint at `https://qrng.anu.edu.au/API/jsonI.php`. Callers should request at least 8 bytes when generating a `uint64` seed.

## Development

- `cd qrand && go test ./...`
