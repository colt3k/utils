# updater Documentation

These notes split the repository-wide docs down to the `updater` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/updater`
- Run `cd updater && go test -mod=mod ./...` in this repository state because vendoring is out of sync.
- Keep shared metadata in `updater.go` aligned with both `artifactory` and `website` backends.
- Treat version comparison, timestamp handling, and download URL construction as high-risk behavior.

## API Surface
- `AppConfig`, `Connection`, and `Version` define the shared metadata types.
- `artifactory.CheckUpdate`, `PerformUpdate`, and `UpdateAvailableMsg` implement the richer host-aware backend.
- `website.CheckUpdate` and `UpdateAvailableMsg` implement the simpler website backend.

## Schema & Data Shapes
- The external update schema is the JSON payload stored in `<app>-<os>-<arch>.update` files.
- The operational file set also includes optional `<app>.auto` flags and `<app>-<os>-<arch>.tgz` archives.

## Dataflow
1. Build local version metadata and one or more update host definitions.
2. Fetch remote update metadata, compare semantic version and build timestamp, and derive the archive URL.
3. Optionally prompt, auto-approve, or download the update depending on the backend and `.auto` state.

## Operator Notes
- This module has a real operator surface because it controls rollout and auto-update behavior.
- If clients stop updating, verify `.update`, `.auto`, archive naming, host reachability, and TLS settings before changing code.
