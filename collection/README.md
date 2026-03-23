# Stack from

[Stack Go](http://github.com/alediaferia/stackgo)

This module keeps a small stack implementation derived from the original `stackgo` project.

## API Surface

- `NewStack()`: create a stack with the default backing size.
- `NewStackWithCapacity(cap int)`: create a stack with a caller-defined block size.
- `Push`, `Pop`, and `Size` provide the expected stack behavior.

## Notes

The implementation is intentionally small and dependency-light. Use it when you want a simple in-memory stack without pulling in a larger collection package.
