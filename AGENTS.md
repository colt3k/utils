# Repository Guidelines

## Project Structure & Module Organization
This repository is a multi-module Go utilities monorepo. Most top-level directories, such as `archive/`, `crypt/`, `netut/`, `store/`, and `timeut/`, are independent modules with their own `go.mod`, tests, and module-level `README.md`. Keep changes scoped to the module you are editing. Subpackages live under their parent module paths, for example `archive/xz`, `io/ioreader/iocsv`, or `netut/hc`. Some modules also vendor dependencies under `vendor/`; only edit vendored code when intentionally updating a vendored dependency.

## Build, Test, and Development Commands
There is no root `go.mod` or `go.work`, so repo-root `go test ./...` does not work. Run commands from the module you changed:

- `cd store && go test ./...` runs that module's test suite.
- `cd netut && go build ./...` builds all packages in the module.
- `cd crypt && go test ./... -run TestLoadCert` runs a focused test.
- `cd mymg && mage -l` lists Mage targets for the build helper module.
- `cd <module> && gofmt -s -w . && go vet ./...` applies standard formatting and a basic static check.

Honor the `go` and `toolchain` versions already declared in each module's `go.mod`.

## Coding Style & Naming Conventions
Follow standard Go formatting and layout. Use tabs as emitted by `gofmt`, keep package names short and lowercase, and use `CamelCase` for exported identifiers. Keep platform-specific behavior in files with Go suffixes such as `*_darwin.go`, `*_linux.go`, and `*_windows.go` instead of mixing OS branches into shared files.

## Testing Guidelines
Place tests next to the code they cover using `*_test.go`. Add regression tests with bug fixes, and prefer table-driven tests for input/output heavy helpers. When changing shared packages, rerun tests in the touched module and any directly affected downstream modules.

## Commit & Pull Request Guidelines
Recent commit subjects are short and imperative, for example `update deps`, `update hash`, and `add new test TestLoadCert`. Follow that pattern: keep the subject brief, describe the behavior change, and mention the affected module when useful. PRs should list touched modules, note OS-specific impact, and include the exact `go test` or `go build` commands you ran.
