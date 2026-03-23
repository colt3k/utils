# io

Provides reader and writer IO methods and abilities.

## Top-Level Helpers

- `reader.go`: open files or URLs, scan lines, read the last lines of a file, and expose CSV helpers.
- `writer.go`: write bytes, strings, temp files, and quoted key/value fields.

## `ioreader`

- `iocsv`
- `ioexif`
- `ioimage`
- `iotab`
- `passthrough` provides ability to display progress or a counter of data

## Other Subpackages

- `iowriter`
- `logstash`

## Development

- `cd io && go test ./...`
