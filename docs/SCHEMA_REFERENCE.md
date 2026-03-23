# Schema Reference

This file remains the repository-level schema overview. Module-specific schema and data-shape notes now live under each subproject’s `docs/README.md`; use [MODULE_INDEX.md](MODULE_INDEX.md) to jump to a module.

## `mymg` Build Configuration (`build.toml`)

`mymg` reads a TOML file shaped like `mymg/test/build.toml`.

```toml
[build]
tags = ""
useAltApps = "yes"

[postclean]
dirs = ["PREP/", "cross"]
files = []

[apps]
gitExe = "/usr/local/bin/git"
tarExe = "/usr/bin/tar"

[[project]]
name = "appname"
ostargets = ["darwin/amd64"]
package = "go.domain.com/colt3k/appname"
version = "cmd/appname/VERSION.txt"
readme = "cmd/appname/README.md"
changelog = "cmd/appname/CHANGES.txt"
files = ["./pkgr/bash_autocomplete"]
```

Important sections:

- `build`: build tags and whether alternate executable paths should be resolved.
- `postclean`: directories and files removed after build/release tasks finish.
- `apps`: absolute paths to helper executables such as `git`, `tar`, `scp`, `sftp`, and hashing tools.
- `scp`, `sftp`: push targets with `host`, `path`, optional `skip_ping`, and `backup`.
- `push-custom`, `pull-custom`: external commands used for custom artifact transport.
- `artifactory`: repository host, path, credentials path, bearer mode, ping behavior, and backup designation.
- `project`: project-specific build metadata including targets, package path, version/readme/changelog files, optional deploy scripts, staged files, and variable overrides.

## Updater Metadata (`updater.AppConfig`)

Remote `.update` files deserialize into `updater.AppConfig`.

```json
{
  "os": "darwin",
  "arch": "arm64",
  "name": "appname",
  "timestamp": 1735689600,
  "version": "v1.2.3",
  "changelog": "https://example.invalid/changelog"
}
```

Operational notes:

- `os`, `arch`, `name`, `timestamp`, `version`, and `changelog` are the externally supplied fields.
- `BaseURL`, `URL`, and `ArchiveName` are derived at runtime by the updater backends.
- `User`, `Pass`, `Bearer`, `DisableVerifyCert`, and `Issue` are transport/runtime fields, not required in the remote JSON.

## Updater Host Configuration (`updater.Connection`)

`updater.Connection` describes a candidate update host. The struct is configured by the embedding application, not by a built-in file parser, but the field meanings are stable:

- Identity and credentials: `Name`, `HostName`, `User`, `PassOrToken`, `Bearer`.
- Repository location: `URLPrefix`, `Repository`, `Path`.
- Reachability gates: `OnAvailable`, `OnAvailableTimeout`, `OnHostNamePrefix`, `OnHostNameSuffix`, `OnAvailableViaHTTP`.
- TLS and backend behavior: `DisableValidateCert`, `AQLSupport`.

## File Metadata Record (`store.FileStore`)

`store.FileStore` is the common serialized record for file inventory work:

- Path metadata: `file`, `name`, `directory`, `symlink`, `file_ext`, `Path`.
- Change tracking: `hash`, `samename`, `status`, `lastmod`, `time`.
- Size metadata: `size`, `sizehr`.
- Media metadata: `exif`, populated lazily for JPEG/JPEG-like image paths.

## Build Metadata Variables (`version`)

The `version` module exposes ldflags-driven variables:

- `VERSION`
- `GITCOMMIT`
- `GITBRANCH`
- `BUILDDATE`
- `GOVERSION`

Typical injection pattern:

```sh
go build -ldflags "-X version.VERSION=v1.2.3 -X version.GITCOMMIT=$(git rev-parse --short HEAD)"
```
