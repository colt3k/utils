# archive Documentation

These notes split the repository-wide docs down to the `archive` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/archive`
- Run `cd archive && go test ./...` from the module root.
- Runnable examples live under `archive/test` and the format-specific `*/test` directories.
- Keep format-specific behavior inside the matching subpackage such as `tgz`, `xz`, or `zip`.

## API Surface
- `archive.Compressor` defines the stream compression interface.
- `tgz.Untar`, `tgz.Tar`, and `tgz.CreateArchive` handle tar.gz creation and extraction.
- `zip.ZipFiles` and `zip.Unzip` handle zip archives.
- `flate.New`, `gz.New`, `lz4.NewLz4`, `lzw.New`, and `xz.NewXz` expose format-specific compressors.

## Schema & Data Shapes
- No external config file schema is defined by this module.
- The main shared shape is a reader/writer pair passed through the `Compressor` interface.

## Dataflow
1. Choose the format-specific package that matches the archive type.
2. Pass source bytes through `Compress` or `Decompress`, or stage files through `tgz`/`zip` helpers.
3. Persist or unpack the resulting stream in the caller.

## Operator Notes
- Library-only module; there is no standing service to operate.
- When changing `tgz` or `zip`, verify behavior on real files in addition to unit tests.
