# store Documentation

These notes split the repository-wide docs down to the `store` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/store`
- Run `cd store && go test ./...` after changing any exported type.
- Keep `FileStore` behavior stable because other modules serialize and sort it.
- Prefer focused container additions over broad new abstractions.

## API Surface
- `bintmap.BiMapInt` provides integer/string bi-directional mapping.
- `FileStore` stores file metadata, status, size, and optional image dimensions.
- `NewMVKeySet`, `StringSet`, and `FormatStore` cover the remaining container helpers.

## Schema & Data Shapes
- No external config file schema is defined by the module.
- `FileStore` is the main reusable record shape and includes path, hash, status, timestamp, size, and EXIF-related fields.

## Dataflow
1. Populate the relevant container or file-metadata type.
2. Mutate or query it through the small exported API.
3. Serialize, sort, or forward the resulting state in higher-level modules.

## Operator Notes
- Library-only module; operational impact comes from downstream consumers that rely on stable metadata fields.
- When changing `FileStore`, check both tests and any serializer or sorter that consumes the struct.
