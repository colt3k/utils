# sembump

Creates version information in the form

major.minor.patch

## Usage

```sh
sembump --kind patch v1.2.3
sembump --kind minor --pre v1.2.3
```

## Notes

The CLI preserves an input `v` prefix and can increment existing prerelease values when `--pre` is used. `mymg` uses this tool during release automation when version bumping is enabled.
