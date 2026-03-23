# Developer Guide

This file remains the repository-level developer overview. Module-specific developer notes are now split under each subproject in `docs/README.md`; start with [MODULE_INDEX.md](MODULE_INDEX.md) when you need the module-local view.

## Repository Model

This repository is a collection of independent Go modules, not a single workspace. Each top-level module owns its own `go.mod`, `go.sum`, README, and tests. Work from the module directory you are changing instead of the repository root.

## Day-To-Day Workflow

1. Pick the target module, for example `store`, `netut`, or `crypt`.
2. Read the module README and nearby tests before changing behavior.
3. Run formatting, tests, and builds from that module:
   - `go test ./...`
   - `go build ./...`
   - `gofmt -s -w .`
4. If the module exposes examples under `test/main.go` or similar, run them when your change affects that flow.

Repo-root `go test ./...` fails because there is no root module. That is expected.

## Code Layout Expectations

- Keep package changes local to the module that owns the import path.
- Prefer platform-specific files such as `*_darwin.go` and `*_windows.go` over runtime OS branching when the behavior is platform-bound.
- Keep tests beside the code they cover using `*_test.go`.
- Preserve existing multi-line comments, URLs, and commented-out code blocks unless you are explicitly extending them.

## Dependencies And Vendoring

Several modules vendor dependencies. Do not edit `vendor/` casually. Update vendored code only when you are intentionally syncing a dependency for that module and can explain the reason in the change.

Modules target mixed Go versions, though many active modules now declare Go `1.24.x` toolchains. Follow the version already pinned in the local `go.mod`.

## Build And Release Tooling

`mymg` is the repository’s build orchestration module. Its Mage targets cover config generation, formatting, linting, testing, install, release, cross-build, and updater flag publication. The reference config lives at `mymg/test/build.toml`.

`sembump` is the version bump CLI used by `mymg` and `gen.go`. When a project config enables version bumping, `mymg` can update version files, README references, git tags, and optional changelog files.

## Documentation Maintenance

- Keep module READMEs aligned with exported APIs and tests.
- Add or update developer/API/schema/dataflow docs when a new operational path or config shape is introduced.
- Prefer short, code-adjacent examples that match the real package names and import paths.
