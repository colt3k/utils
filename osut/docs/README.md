# osut Documentation

These notes split the repository-wide docs down to the `osut` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/osut`
- Run `cd osut && go test ./...` after changing platform-specific behavior.
- Keep OS-specific code in `*_darwin.go`, `*_linux.go`, and `*_windows.go` files.
- Do not merge platform branching back into shared files unless the logic is truly cross-platform.

## API Surface
- `OS`, `Windows`, `Linux`, `Mac`, and `Android` expose platform detection.
- `Hostname`, `OSVersionMaj`, `OSVersionMinor`, `OSDistro`, and `FullVer` expose host metadata.
- `CallCmd`, `CallCmdNoWait`, and `FindProcess` expose command/process helpers.
- `RunAppWithPath`, `RunAppNoPath`, and `RunAppMac` launch applications.

## Schema & Data Shapes
- No external schema file is used.
- Primary data shapes are `Platform`, `Shell`, and the OS-specific version data returned by the helper functions.

## Dataflow
1. Detect the current platform and version information.
2. Call the appropriate command or executable helper for the OS.
3. Return normalized platform data or process results to the caller.

## Operator Notes
- Library module only, but it executes real OS commands.
- Test platform-specific changes on the matching OS instead of assuming behavior from one platform transfers to another.
