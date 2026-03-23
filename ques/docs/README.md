# ques Documentation

These notes split the repository-wide docs down to the `ques` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/ques`
- Run `cd ques && go test ./...` when adding tests or changing prompt behavior.
- Keep the package small and blocking; higher-level menu orchestration belongs in `imnu` or caller code.
- Avoid hidden side effects beyond reading input and returning values.

## API Surface
- `Question`, `QuestionOpts`, `QuestionInt`, and `Confirm` cover the supported prompt styles.

## Schema & Data Shapes
- No external schema is used.
- Inputs are prompt text and optional string options; outputs are scalar values.

## Dataflow
1. Choose the prompt helper that matches the expected answer type.
2. Display the prompt and read from stdin.
3. Return the parsed value to the caller.

## Operator Notes
- CLI-only helper; there is no service runbook.
- Verify prompt wording and parsing behavior when integrating with scripts or release tooling.
