# version

Simple reusable data fields between applications.

## Exported Variables

- `VERSION`
- `GITCOMMIT`
- `GITBRANCH`
- `BUILDDATE`
- `GOVERSION`

## Usage

Populate these values with `-ldflags` during `go build` so applications can report the exact release and build metadata they are running.
