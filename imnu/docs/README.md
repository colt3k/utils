# imnu Documentation

These notes split the repository-wide docs down to the `imnu` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/imnu`
- Run `cd imnu && go test ./...` when adding tests or changing selection behavior.
- Keep menu logic simple and caller-driven; the caller owns enabled-state checks and task side effects.
- Use the README example as the baseline interaction model.

## API Surface
- `Menu` defines the id, label, description, enable predicate, and task callback.
- `InteractiveMenu` manages the menu loop.
- `New`, `CaptureSelection`, and `Pause` are the primary exported entry points.

## Schema & Data Shapes
- No external schema is used.
- The main runtime shapes are `Menu`, `InteractiveMenu`, and the selection map passed to `CaptureSelection`.

## Dataflow
1. Build a slice of `Menu` definitions with enable predicates and tasks.
2. Create an `InteractiveMenu` with `New` and start the loop.
3. Let each selected task mutate caller state before the next menu render.

## Operator Notes
- CLI helper module only; there is no long-running operator surface.
- When changing prompts, verify that disabled entries remain hidden or unselectable as intended.
