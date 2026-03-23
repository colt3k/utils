# timeut Documentation

These notes split the repository-wide docs down to the `timeut` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/timeut`
- Run `cd timeut && go test ./...` after changing date or calendar logic.
- Be careful with time zones and week calculations because they are easy to regress subtly.
- Prefer additive helpers over hidden global time behavior.

## API Surface
- `ParseRFC3339`, `ConvertUnix2Time`, `ConvertUnix2TimeStr`, and `ConverMillis2Time` cover parsing and conversion.
- `Time(...)` and `GMTTime()` construct the `MyTime` helper.
- `StartTime`, `StartDate`, `Julian2Date`, and `Date2Julian` cover calendar calculations.

## Schema & Data Shapes
- No external schema is used.
- The primary data shapes are `time.Time`, `time.Duration`, and the module's `MyTime` wrapper.

## Dataflow
1. Parse or construct the relevant time value.
2. Apply the conversion or calendar helper that matches the use case.
3. Return the derived time, date components, or Julian day value to the caller.

## Operator Notes
- Library-only module with no service runbook.
- If date math looks wrong, verify timezone assumptions and week/year boundaries before changing code.
