# mymg Documentation

These notes split the repository-wide docs down to the `mymg` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/mymg`
- Use `cd mymg && mage -l` to inspect targets and `cd mymg/test && mage display` to inspect the sample config.
- Keep `build.toml` parsing logic aligned with the sample config in `mymg/test/build.toml`.
- Do not change deploy-script semantics casually because packaging, auto-update publication, and release automation depend on them.

## API Surface
- Main operator-facing targets include `GenConf`, `Display`, `Build`, `Install`, `Release`, `BuildCross`, `Auto`, `AutoStatus`, `NoAuto`, `Format`, `Test`, and `Vet`.
- Core config shapes include `Config`, `Project`, `Apps`, `BuildData`, `PostClean`, `ScpData`, `SftpData`, `PushCustom`, `PullCustom`, and `ArtifactoryData`.

## Schema & Data Shapes
- The external schema is the TOML file shaped like `mymg/test/build.toml`.
- Sections include `[build]`, `[postclean]`, `[apps]`, transport arrays such as `[[scp]]` and `[[artifactory]]`, and one or more `[[project]]` entries.

## Dataflow
1. Parse `build.toml` into the config structs and select enabled projects.
2. Run format, lint, test, vet, build, package, and optional version-bump steps.
3. Publish artifacts through SCP, SFTP, custom commands, or Artifactory, and manage `.auto` files for the updater.

## Operator Notes
- This module has a real operator surface: verify helper executable paths, credentials, packaging files, and remote destinations before release.
- Use the sample project under `mymg/test` as the canonical dry-run environment before touching production build configs.
