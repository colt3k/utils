# store

Various data storage objects

- bimap
- file store
- multi value keyset
- string set

## Key Types

- `bintmap.BiMapInt`: integer-to-string bi-directional lookup
- `FileStore`: serialized file metadata and optional image dimensions
- `MVKeySetMap`: one key to many values
- `StringSet`: map-backed string membership checks
- `FormatStore`: generic `[]interface{}` transport helper

## Development

- `cd store && go test ./...`
