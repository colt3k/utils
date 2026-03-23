# Dataflow

This file captures repository-level flows that cross module boundaries. Module-specific dataflow notes now live under each subproject’s `docs/README.md`; use [MODULE_INDEX.md](MODULE_INDEX.md) for the split view.

## Build And Release Flow (`mymg`)

1. `mymg` parses `build.toml` into `Config`, `Projects`, `Apps`, and transport sections.
2. Selected projects run `Format`, `Lint`, `Test`, and `Vet` before a build artifact is produced.
3. `Build` compiles the project, applies ldflags, and optionally compresses the binary with UPX.
4. `Release` and `BuildCross` stage extra files into `PREP/`, build target archives, and publish through SCP, SFTP, custom commands, or Artifactory.
5. `Auto`, `NoAuto`, and `AutoStatus` manage the `.auto` sidecar used by the updater flow.

## Runtime Config Lookup Flow (`config`)

`config.Config.Load` resolves the requested config file in three passes:

1. executable directory
2. caller file directory
3. current working directory

Once found, `go-up` loads the file into the `Config.Util` handle. Logging is delegated to the configured `Logger`.

## Update Check And Download Flow (`updater`)

1. The application builds a local `updater.Version` and one or more `updater.Connection` values.
2. The backend checks host reachability and selects the first viable remote.
3. It downloads `<app>-<os>-<arch>.update` and optionally `<app>.auto`.
4. Remote metadata is decoded into `updater.AppConfig`.
5. Version comparison is done by semantic version first, then build timestamp as a fallback.
6. When newer content is found, the backend derives the archive URL `<app>-<os>-<arch>.tgz`.
7. `PerformUpdate` either reports availability, prompts the user, or downloads automatically when `.auto` is set to `true`.

## File Inventory And Sorting Flow

1. `store.FileStore` records are populated with path, size, hash, and modification metadata.
2. `GetEXIF` lazily enriches JPEG records by reading image dimensions through `io/ioreader/ioimage`.
3. `sort.Stores` provides case-insensitive ordering for `[]store.FileStore`.
4. File and IO helpers from `file` and `io` handle copy, rotation, reading, writing, and pass-through progress reporting.

## HTTP And Network Request Flow (`netut/hc`)

1. Callers compose a client through `hc.NewClient` or the higher-level `processor.MakeCall`.
2. Optional auth, TLS settings, timeouts, and redirect behavior are applied through `ClientOption` or `HTTPClientSettings`.
3. Requests return parsed response data, status code, and optional tracing information for diagnostics.
