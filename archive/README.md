# Archiving utilities

Stream and file archiving helpers shared across the repository. The root package exposes the common `archive.Compressor` interface, while format-specific subpackages provide concrete implementations.

## Formats

- `flate`
- `gz`
- `lz4`
- `lzw`
- `tgz`
- `xz`
- `zip`

## Package Map

- `archive/tgz`: create or extract tar.gz content from directories or explicit file lists.
- `archive/zip`: zip and unzip files and folders.
- `archive/gz`, `archive/lz4`, `archive/lzw`, `archive/xz`, `archive/flate`: stream-oriented compressors that fit the shared compressor pattern.

## Examples And Tests

- `go run ./test`
- `go run ./xz/test`
- `cd archive && go test ./...`
