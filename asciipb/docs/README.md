# asciipb Documentation

These notes split the repository-wide docs down to the `asciipb` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/ascii`
- Run `cd asciipb && go test ./...` for module checks.
- Use `go run ./test` to view the progress output interactively.
- Keep terminal escape handling in `CursorNav` rather than scattering raw escape strings.

## API Surface
- `ProgressIndicator(size int)` renders a simple progress indicator.
- `ProgressBarIndicator(size int)` renders a single bar.
- `MultiProgressBarIndicator(size int)` handles repeated progress output.
- `CursorNav` stores the cursor movement tokens used by the package.

## Schema & Data Shapes
- No external config or file schema is used.
- The only input shape is the integer `size` passed to the rendering functions.

## Dataflow
1. Select the indicator style that matches the CLI task.
2. Pass the work size into the renderer while the task is running.
3. Let the caller own task progress and output cadence.

## Operator Notes
- Library-only module with no operator-managed runtime.
- Check output on a real terminal before changing escape behavior.
