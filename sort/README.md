# sort

Home to sorting utilities.

## Current Scope

The module currently exposes `Stores`, a sorter for `[]store.FileStore` values. Ordering is case-insensitive first, with the original rune case used as a tie-breaker.

## Notes

The import path is `github.com/colt3k/utils/sort`, but the package name inside `sorter.go` is `utils`. Alias it explicitly if that improves readability in callers.
