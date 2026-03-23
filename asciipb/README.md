# asciipb - ASCII Progressbar

Small terminal progress indicators for command-line tools.

## API Surface

- `ProgressIndicator(size int)`: simple spinner-like progress output.
- `ProgressBarIndicator(size int)`: single progress bar.
- `MultiProgressBarIndicator(size int)`: progress output for repeated operations.

## Notes

`CursorNav` holds the cursor navigation escape sequences used by the package. For a runnable example, see `asciipb/test/main.go`.

## Development

- `cd asciipb && go test ./...`
- `go run ./test`
