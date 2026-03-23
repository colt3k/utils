# logstash Documentation

These notes split the repository-wide docs down to the `logstash` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/logstash`
- Run `cd logstash && go test -mod=mod ./...` in this repository state because vendoring is out of sync.
- Keep connection lifecycle and deadline refresh behavior explicit; downstream code treats this as a network transport.
- Verify both plain TCP and TLS paths when adjusting connection logic.

## API Surface
- `Server` stores connection settings and the active socket.
- `New`, `Connect`, and `ConnectTLS` create the client connection.
- `Write`, `WriteTLS`, `Close`, and `CloseTLS` manage event delivery and connection teardown.

## Schema & Data Shapes
- No config file schema is defined inside the module.
- TLS mode expects a CA PEM file at `~/mycerts/lsCA.pem` and uses the configured host as the TLS server name.

## Dataflow
1. Construct a `Server` with host, port, and timeout.
2. Open a plain or TLS connection depending on the target.
3. Write newline-delimited payloads and refresh deadlines after each successful send.

## Operator Notes
- This is not a daemon, but it depends on live network and certificate state.
- If TLS delivery fails, confirm the CA file path, server name, and timeout settings before changing code.
