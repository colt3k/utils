# collection Documentation

These notes split the repository-wide docs down to the `collection` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/collection`
- Run `cd collection && go test ./...` when adding tests or changing behavior.
- The code is derived from `stackgo`; preserve the lightweight API shape.
- Prefer focused additions over turning this module into a general-purpose container library.

## API Surface
- `Stack` is the exported stack type.
- `NewStack()` creates a stack with the default backing capacity.
- `NewStackWithCapacity(cap int)` creates a stack with a caller-defined capacity.
- `Push`, `Pop`, and `Size` provide the main operations.

## Schema & Data Shapes
- No external schema is defined.
- The module only exposes in-memory stack state.

## Dataflow
1. Create a stack instance.
2. Push elements as work accumulates.
3. Pop values in LIFO order until the stack is empty.

## Operator Notes
- Library-only module; there is no runbook beyond normal tests.
- If you change allocation behavior, add or update regression coverage first.
