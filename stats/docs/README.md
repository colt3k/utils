# stats Documentation

These notes split the repository-wide docs down to the `stats` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/stats`
- Run `cd stats && go test ./...` after changing numerical behavior.
- Preserve upstream numerical expectations unless you are intentionally diverging and updating tests.
- Use the existing distribution and test coverage as the baseline for compatibility.

## API Surface
- `Sample`, `Bounds`, `Mean`, `Variance`, and `StdDev` cover descriptive statistics.
- `NormalDist`, `TDist`, `UDist`, and `DeltaDist` expose distribution types.
- `TwoSampleTTest`, `PairedTTest`, `OneSampleTTest`, and `MannWhitneyUTest` expose hypothesis tests.
- `InvCDF` and `Rand` provide helper constructors for distributions.

## Schema & Data Shapes
- No external config schema is used.
- Core data shapes are float slices, distribution structs, and typed test result structs such as `TTestResult` and `MannWhitneyUTestResult`.

## Dataflow
1. Build sample slices or distribution instances.
2. Run the descriptive statistic or hypothesis test API that matches the analysis.
3. Interpret the returned result struct or numeric value in the caller.

## Operator Notes
- Library-only module with no operator-managed runtime.
- Numerical regressions can be subtle, so compare against tests and known-good statistical references before merging changes.
