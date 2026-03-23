# mathut Documentation

These notes split the repository-wide docs down to the `mathut` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/mathut`
- Run `cd mathut && go test ./...` after changing parsing or percentile logic.
- Keep helpers deterministic and side-effect free.
- Review downstream modules such as `updater` and `osut` when formatting behavior changes.

## API Surface
- `Round`, `FmtFloat*`, and `FmtInt` cover formatting.
- `ParseFloat`, `ParseInt`, and `IntSize` cover parsing and sizing.
- `Percentile`, `Median`, and `PercentDiff` cover common numeric calculations.

## Schema & Data Shapes
- No external schema is used.
- Inputs are primitive numeric or string values, and outputs are primitive values or errors.

## Dataflow
1. Normalize or parse the caller's numeric input.
2. Run the desired calculation or formatter.
3. Return the derived value to the calling package.

## Operator Notes
- Library-only module with no service runbook.
- Small rounding or parsing changes can cascade into many modules, so validate downstream formatting expectations.
