# Configuration Object

Small configuration loader with pluggable logging.

## Current Behavior

- provides `Load`, `Save`, and `Delete` functions
- `Save` and `Delete` are still placeholders
- `Load` searches for the requested file in the executable directory, the caller directory, and finally the current working directory

## Logger Support

The module exposes a small `Logger` interface plus `Standard` and `Nop` implementations so applications can control how configuration events are emitted.

## Development

- `cd config && go test ./...`
- keep README and code comments aligned if the config lookup order changes
