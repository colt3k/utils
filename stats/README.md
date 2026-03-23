## Exported library from Golang internal source

https://pkg.go.dev/golang.org/x/perf/internal/stats

https://github.com/golang/perf

This package was exported for use by other libraries and projects, due to it being internal only on original project.

## Included Functionality

- distributions such as normal, t, delta, and uniform
- descriptive statistics through `Sample`, `Bounds`, `Mean`, `Variance`, and `StdDev`
- hypothesis tests including `TwoSampleTTest`, `PairedTTest`, `OneSampleTTest`, and `MannWhitneyUTest`

## Notes

The implementation is a trimmed fork and keeps the upstream license in this module. Use the local tests as the primary compatibility check when changing behavior.

## Development

- `cd stats && go test ./...`
