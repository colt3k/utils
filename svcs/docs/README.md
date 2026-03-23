# svcs Documentation

These notes split the repository-wide docs down to the `svcs` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/svcs`
- Run `cd svcs && go test -mod=mod ./...` in this repository state because vendoring is out of sync.
- Keep the generated unit file and `/opt/gosvc` path assumptions explicit; operators depend on them.
- Treat all `sudo` and `systemctl` shell-outs as high-impact behavior.

## API Surface
- The CLI accepts `-service`, `-runas`, and `-remove` flags.
- `ServiceTemplate` is the generated systemd unit file body.
- `Sudo`, `SudoNoFail`, `Command`, and `PanicOnError` support the main workflow.

## Schema & Data Shapes
- There is no config file schema.
- The key runtime shapes are CLI flags and the generated systemd unit written under `/opt/gosvc`.

## Dataflow
1. Build the service binary from `<service>/<service>.go`.
2. Write the systemd unit file into `/opt/gosvc` and enable it with `systemctl`.
3. Start, stop, disable, or remove the service depending on the selected flags.

## Operator Notes
- This module has a real operator surface because it modifies system services.
- Use it only on hosts where `/opt/gosvc`, `sudo`, and `systemctl` are expected and controlled.
