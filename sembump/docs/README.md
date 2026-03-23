# sembump Documentation

These notes split the repository-wide docs down to the `sembump` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/sembump`
- Run `cd sembump && go test ./...` if you add tests or modify version parsing.
- Keep CLI output machine-friendly because other tooling, especially `mymg`, consumes it.
- Preserve `v` prefix handling and prerelease increment semantics.

## API Surface
- The CLI accepts `--kind` / `-k` with `major`, `minor`, or `patch`.
- `--pre` creates or increments prerelease variants.
- Input and output version strings preserve a leading `v` when present.

## Schema & Data Shapes
- No config file schema is used.
- The only required input is a semantic version string supplied on the command line.

## Dataflow
1. Parse flags and the input version string.
2. Normalize the version into semantic-version components.
3. Bump the requested segment and print the new version to stdout.

## Operator Notes
- Operator-facing only as a CLI utility, typically inside build or release automation.
- If version bumps look wrong, verify the input string and `--pre` mode before inspecting the code.
