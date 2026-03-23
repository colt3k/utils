# msngr Documentation

These notes split the repository-wide docs down to the `msngr` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/msngr`
- Run `cd msngr && go test ./...` after changing any backend.
- Keep platform-specific implementations in their matching subpackages such as `freedesktop`, `nsuser`, or `speechsynthesizer`.
- Treat provider-specific request shapes as backend concerns rather than forcing a single global message schema.

## API Surface
- Each backend exposes its own `Notification` type and send logic inside its subpackage.
- Supported integrations include `espeak`, `freedesktop`, `notifyicon`, `nsuser`, `pushbullet`, `say`, `slack`, and Windows speech synthesis.

## Schema & Data Shapes
- There is no single external config schema in this module.
- Each backend owns its own request shape because OS and provider requirements differ.

## Dataflow
1. Choose the backend that matches the target OS or provider.
2. Build the backend-specific notification payload.
3. Dispatch the message through the selected transport.

## Operator Notes
- Operator steps depend on the backend: local desktop permissions, provider tokens, and OS support all matter.
- When incidents are backend-specific, test that backend directly rather than assuming parity across platforms.
