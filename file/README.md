# File

Filesystem helpers and common file-related interfaces used by other modules.

## Provides

- native file utils
- file permissions
- mimetypes
- sizes
- file meta data

## Key Areas

- `utils.go`: copy, delete, rotate, and path expansion helpers.
- `filenative`: concrete `file.File` implementation for local files.
- `filemeta`: metadata extraction, with platform-specific implementations where needed.
- `fileperm`: permission bit inspection and modification.
- `mimetypes` and `filesize`: MIME lookup and size conversion helpers.

## Development

- `cd file && go test ./...`
