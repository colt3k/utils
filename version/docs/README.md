# version Documentation

These notes split the repository-wide docs down to the `version` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/version`
- Run `cd version && go test ./...` if you add tests or change the exported variables.
- Keep variable names stable because build scripts and applications set them with `-ldflags`.
- Document any new variable in module README and cross-repo docs before using it widely.

## API Surface
- `VERSION`, `GITCOMMIT`, `GITBRANCH`, `BUILDDATE`, and `GOVERSION` are the exported build metadata variables.

## Schema & Data Shapes
- There is no config file schema.
- The effective input schema is the linker flag set used during `go build`.

## Dataflow
1. Inject values with `-ldflags` during the build.
2. Import the module where the application needs version metadata.
3. Report or serialize the build metadata at runtime.

## Operator Notes
- Library-only module, but it participates directly in release metadata.
- If version output is wrong, inspect the build command and linker flags before changing code.
