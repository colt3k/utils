# Operator Runbook

This file is the repository-level operator view. Module-local operator guidance now lives under each subproject’s `docs/README.md`; use [MODULE_INDEX.md](MODULE_INDEX.md) for the split docs.

## Preconditions

- Use the Go version pinned by the target module’s `go.mod`.
- Run commands from the module or application directory, not the repository root.
- For `mymg` releases, confirm that external tools referenced in `build.toml` exist: `git`, `tar`, `scp`, `sftp`, hash utilities, optional `upx`, and any custom push/pull commands.
- Make sure release credentials and Artifactory token files are present outside version control.

## Build And Verify Locally

1. Validate build config:
   - `cd mymg/test`
   - `mage -l`
   - `mage display`
2. Build the current target:
   - `mage build`
   - `mage install`
3. Verify staged or installed artifacts:
   - confirm the binary starts
   - confirm any README, VERSION, and packaged support files are present
   - run the target module’s `go test ./...`

## Release And Publish

1. Use `mage release` for the standard release path.
2. Use `mage buildcross` when cross-platform archives are required.
3. If version bumping is enabled, `mymg` can update version files, README references, git tags, and optional changelog files before publication.
4. Review the generated archive names and remote paths before pushing to SCP, SFTP, custom transport, or Artifactory targets.

## Auto-Update Operations

- `mage auto`: publish `.auto` with `true` to enable unattended updates.
- `mage noauto`: publish `.auto` with `false` to require manual confirmation.
- `mage autostatus`: fetch and verify the currently published auto-update flag.

The updater expects `.update`, `.auto`, and `<app>-<os>-<arch>.tgz` artifacts to stay consistent.

## Rollback

1. Disable unattended rollout with `mage noauto`.
2. Publish the last known-good archive and matching `.update` metadata.
3. Verify both semantic version and timestamp values. The updater falls back to timestamp comparison when versions do not move forward, so a rollback must not accidentally look newer than the intended release.
4. Re-run `mage autostatus` and smoke-test an update check from a client build.

## Troubleshooting

- If update checks fail early, verify host reachability and TLS settings in `updater.Connection`.
- If `mage` cannot find helper executables, fix the `apps` section in `build.toml`.
- If UPX is missing or unsupported, `mymg` will skip binary compression.
- If tests fail at repository root, rerun from the target module because this repo is not a single Go workspace.
