# mymg

Overlay to Mage that provides a config-driven framework to build and deploy Go projects.

## Core Targets

- `mage genconf`: generate a starter `build.toml`
- `mage display`: parse and display the active config
- `mage build`: local build with format, lint, test, and vet steps
- `mage install`: install the selected project locally
- `mage release`: build, package, and publish release artifacts
- `mage buildcross`: produce cross-platform artifacts
- `mage auto`, `mage noauto`, `mage autostatus`: manage updater auto-rollout flags

## Configuration

The reference schema lives in `mymg/test/build.toml`. Important sections are `build`, `postclean`, `apps`, transport targets, and one or more `[[project]]` entries.

## Notes

`mymg` can bump versions, update README references, create git tags, and publish `.auto` files for the updater module. The sample project under `mymg/test` is the best place to see the expected file layout.
