# sort Documentation

These notes split the repository-wide docs down to the `sort` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/sort`
- Run `cd sort && go test -mod=mod ./...` in this repository state because vendoring is out of sync.
- Note that the module import path is `github.com/colt3k/utils/sort` but the package name in code is `utils`.
- Preserve the current case-insensitive comparison semantics unless consumers are updated too.

## API Surface
- `Stores` implements `Len`, `Swap`, and `Less` for `[]store.FileStore`.
- `Less` compares names case-insensitively and uses original rune case as a tie-breaker.

## Schema & Data Shapes
- No external schema file is defined.
- The effective input shape is a slice of `store.FileStore` values.

## Dataflow
1. Build a slice of `store.FileStore` values.
2. Wrap it as `sort.Stores` in caller code.
3. Use the standard library sorting helpers to order the slice.

## Operator Notes
- Library-only module with no runtime operator steps.
- If ordering changes unexpectedly, verify mixed-case filenames and name-prefix edge cases first.
