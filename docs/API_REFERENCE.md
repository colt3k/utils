# API Reference

This document is a repository index for the main public APIs. Module-specific API notes now live under each subproject’s `docs/README.md`; use [MODULE_INDEX.md](MODULE_INDEX.md) for direct links.

## Build, CLI, And Operator Packages

- `mymg`: Mage targets such as `GenConf`, `Display`, `Build`, `Release`, `BuildCross`, `Auto`, `NoAuto`, `Format`, `Test`, and `Vet`.
- `sembump`: semantic version CLI that bumps `major`, `minor`, or `patch`, with optional prerelease handling.
- `svcs`: systemd installer/remover CLI for building and registering a Go service binary.
- `imnu`: terminal menus via `New`, `CaptureSelection`, and `Pause`.
- `ques`: interactive prompts via `Question`, `QuestionOpts`, `QuestionInt`, and `Confirm`.
- `asciipb`: terminal progress indicators via `ProgressIndicator`, `ProgressBarIndicator`, and `MultiProgressBarIndicator`.

## Runtime And Process Helpers

- `config`: `Config`, `Logger`, `Standard`, and `Nop` for environment-backed configuration loading.
- `env`: environment lookup and filtering via `Find`, `Prefix`, `Suffix`, `Includes`, `All`, and `Add`.
- `debug`: pprof snapshot helpers such as `Stack`, `Heap`, `ThreadCreate`, and `Block`.
- `shutdown`: singleton signal handling through `SetupNotifyContext` and `Graceful`.
- `ctxt`: `SleepContext` for cancellation-aware delays.
- `lock`: file-lock coordination via `New`, `Try`, and `Unlock`.
- `osut`: OS detection, command execution, version lookup, and `runner` helpers.
- `profile`: lightweight runtime metrics via `Duration`, `MemUsage`, and `CpuUsagePercent`.
- `version`: build metadata variables for version, branch, commit, build date, and Go version.

## Storage, Files, And Data Movement

- `archive`: `Compressor` interface plus concrete implementations in `flate`, `gz`, `lz4`, `lzw`, `tgz`, `xz`, and `zip`.
- `file`: filesystem helpers, file metadata, permissions, MIME lookup, size formatting, and rotation helpers.
- `io`: generic readers/writers plus `iocsv`, `ioexif`, `ioimage`, `iotab`, `passthrough`, and logstash writers.
- `sqlite`: singleton database wrapper via `DB`, `Execute`, `Query`, and `Close`.
- `store`: `FileStore`, string sets, bi-maps, generic stores, and multi-value key sets.
- `sort`: `Stores` sorter for `[]store.FileStore`.
- `print`: table output through `Printer.TablePrint`.
- `logstash`: TCP/TLS writer for sending log lines to a Logstash-compatible endpoint.

## Crypto, Network, And Delivery

- `crypt`: key derivation, AES helpers, streaming encryption, Poly1305 signing, keypair generation, certificate helpers, and signature verifiers.
- `hash`: shared `Hasher` API, HMAC helpers, and algorithm-specific constructors for MD, SHA, and BLAKE2 variants.
- `netut`: IP helpers, TCP/UDP clients and servers, HTTPS server setup, certificate utilities, basic auth middleware, and configurable HTTP clients under `hc`.
- `msngr`: OS-specific notification integrations for `espeak`, `freedesktop`, `nsuser`, `notifyicon`, `pushbullet`, `say`, `slack`, and Windows speech synthesis.
- `updater`: shared update metadata types plus `artifactory` and `website` backends for checking, advertising, and downloading updates.
- `webut`: JSON helper `ConvertibleBoolean` for booleans encoded as strings.
- `qrand`: QRNG-backed seed generation using the ANU quantum random number service.

## Utility Packages

- `concur`: bounded worker pools and a resynchronizing `Once`.
- `retry`: exponential/backoff retry processing with context support.
- `repeat`: fixed-interval repeating tasks with cancellation and stop hooks.
- `mathut`: rounding, formatting, parsing, percent-diff, percentiles, and median helpers.
- `stats`: descriptive statistics, distributions, and t-test / U-test implementations.
- `stringut`: string conversion, extraction, validation, and normalization helpers.
- `timeut`: RFC3339 parsing, Unix and millisecond conversion, GMT helpers, week starts, and Julian date conversion.
- `collection`: stack implementation derived from `stackgo`.
- `img`: RGB image decoding helper.
- `regx`: common regular-expression constants and validators.
