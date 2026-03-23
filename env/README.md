# Environment

Loads and works with environment variables.

## API Surface

- `Find(key)` returns a trimmed environment value.
- `Prefix`, `Suffix`, and `Includes` filter the environment map by key pattern.
- `All()` returns a snapshot map of the current environment.
- `Add(key, val)` writes a variable into the current process environment.

## Notes

Top-level helpers delegate to a shared `Environment` instance, so callers can use either the package functions or the struct methods.

## Development

- `cd env && go test ./...`
