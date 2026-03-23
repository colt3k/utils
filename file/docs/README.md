# file Documentation

These notes split the repository-wide docs down to the `file` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/file`
- Run `cd file && go test ./...` after changing file operations or platform-specific code.
- Keep OS-specific logic in `filemeta/*` or `fileperm/*` rather than branching heavily in shared code.
- Treat path handling and deletion helpers as sensitive because downstream modules depend on them heavily.

## API Surface
- `utils.go` provides copy, delete, path expansion, and file rotation helpers.
- `file.File`, `Meta`, `Exif`, and `Image` define the common file interfaces.
- `filenative.NewFile` creates a native file implementation.
- `filemeta.New`, `fileperm.Stat`, and `mimetypes.Find` cover metadata, permission, and MIME lookups.

## Schema & Data Shapes
- No external config file schema is defined by this module.
- Primary data shapes are path strings, `os.FileMode`-derived permission bits, and file metadata structs.

## Dataflow
1. Resolve or normalize the file path.
2. Use the appropriate helper for copy, delete, rotation, metadata, or permission work.
3. Return lightweight interfaces so higher-level modules can compose behavior without depending on OS details.

## Operator Notes
- Library-only module; failures usually surface as path, permission, or cleanup bugs in downstream code.
- When changing destructive helpers such as delete or rotate functions, verify on disposable test data first.
