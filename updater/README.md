# Updater

Provides an ability to perform an application update from either Artifactory or a typical website.
 
Start HTTP for testing
# Python 3
python -m http.server 8081 

# Python 2
python -m SimpleHTTPServer 8081

## Included Packages

- `updater`: shared metadata and connection types
- `artifactory`: host-aware update checks, changelog fetching, and archive download
- `website`: simple HTTP-hosted update checks

## Runtime Files

The updater flow expects:

- `<app>-<os>-<arch>.update`
- optional `<app>.auto`
- `<app>-<os>-<arch>.tgz`

## Notes

Version comparison prefers semantic version and falls back to build timestamp when needed. See `docs/SCHEMA_REFERENCE.md` and `docs/DATAFLOW.md` for the expected metadata format and runtime sequence.

## Development

- `cd updater && go test ./...`
