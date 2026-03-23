# Hash

Provides easy utilities for hashing.

## Algorithms

- MD variants
- SHA-1
- SHA-2
- SHA-3
- SHA-512
- BLAKE2
- HMAC helpers

## API Surface

Use `hash.String` or `hash.File` with a standard `hash.Hash`, or choose a package-specific constructor such as `sha2.NewHash(...)` or `blake2.NewHash(...)` when you want formatted output helpers.

## Development

- `cd hash && go test ./...`
