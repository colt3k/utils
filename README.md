# colt3k/utils

`github.com/colt3k/utils` is a multi-module Go utility repository. Most top-level directories are independent modules with their own `go.mod`, tests, examples, and release cadence. The repo combines build tooling, CLI helpers, storage and IO packages, network helpers, crypto primitives, and a small updater stack.

## Repository Layout

- `archive`, `file`, `hash`, `io`, `sqlite`, `store`: data movement, metadata, and storage helpers.
- `crypt`, `netut`, `osut`, `msngr`, `updater`: security, networking, operating-system, notification, and update delivery helpers.
- `concur`, `retry`, `repeat`, `shutdown`, `stats`, `stringut`, `timeut`, `mathut`: reusable runtime and utility packages.
- `mymg`, `sembump`, `svcs`: build, release, and service-management tooling.

## Working In This Repo

There is no root `go.mod` or `go.work`, so work inside the module you are changing.

- `cd <module> && go test ./...`
- `cd <module> && go build ./...`
- `cd mymg && mage -l`

Many modules vendor dependencies. Treat `vendor/` as mirrored dependency state unless you are intentionally updating a vendored package.

## Documentation Map

- [AGENTS.md](AGENTS.md): contributor workflow and repo-specific guardrails.
- [docs/MODULE_INDEX.md](docs/MODULE_INDEX.md): links to the per-module technical docs split under each subproject.
- [docs/DEVELOPER_GUIDE.md](docs/DEVELOPER_GUIDE.md): repo-level development workflow and testing expectations.
- [docs/API_REFERENCE.md](docs/API_REFERENCE.md): repo-level package index and notable exported APIs.
- [docs/SCHEMA_REFERENCE.md](docs/SCHEMA_REFERENCE.md): repo-level shared schemas such as `build.toml` and updater metadata.
- [docs/DATAFLOW.md](docs/DATAFLOW.md): repo-level build, config, storage, and update execution flow.
- [docs/OPERATOR_RUNBOOK.md](docs/OPERATOR_RUNBOOK.md): repo-level release, rollout, verification, and rollback procedures.

## Module Docs

Each module README documents its public entry points, examples, and any platform-specific behavior. Each module now also has a `docs/README.md` with developer notes, API surface, schema/data-shape notes, dataflow, and operator guidance. Packages with OS-specific implementations follow standard Go suffixes such as `*_darwin.go`, `*_linux.go`, and `*_windows.go`.
