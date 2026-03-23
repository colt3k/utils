# env Documentation

These notes split the repository-wide docs down to the `env` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/env`
- Run `cd env && go test ./...` after changing filter or case-sensitivity logic.
- Keep top-level helpers and `Environment` methods behaviorally aligned.
- Avoid adding OS-specific behavior here; callers should normalize separately.

## API Surface
- `Find`, `Prefix`, `Suffix`, `Includes`, `All`, and `Add` expose package-level accessors.
- `Environment` provides the same operations as instance methods for custom callers.

## Schema & Data Shapes
- This module works against the current process environment.
- Returned data shapes are plain `string` values or `map[string]string` snapshots.

## Dataflow
1. Read a single key or build a filtered view using prefix, suffix, or substring matching.
2. Adjust case handling with the `caseSensitive` flag.
3. Write new variables into the process environment with `Add` when needed.

## Operator Notes
- Library-only module with no standalone runbook.
- If a filter returns unexpected results, inspect the caller's environment and the `caseSensitive` flag first.
