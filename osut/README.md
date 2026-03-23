# OS Utilities

- OS Version
- Name
- Platform
- OS Shell

Mac, Linux, Windows and Android.

## Key Areas

- `os.go`: platform, hostname, and distribution helpers
- `oscmds.go`: command execution and process lookup
- `runner.go`: launch apps with or without an explicit path
- `osexec`: resolve executable and executable folder locations

## Platform Notes

OS version detection is split into `*_darwin.go`, `*_linux.go`, and `*_windows.go` files. Keep new platform-specific behavior in the matching file instead of merging it into the shared package path.

## Development

- `cd osut && go test ./...`
