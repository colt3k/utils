# logstash

TCP and TLS client helpers for writing newline-delimited messages to a Logstash-compatible endpoint.

## API Surface

- `New(host, port, timeout)`
- `Connect`, `ConnectTLS`
- `Write`, `WriteTLS`
- `Close`, `CloseTLS`

## Notes

TLS mode expects a CA file at `~/mycerts/lsCA.pem`. The server object keeps its own timeout configuration and refreshes deadlines after each successful write.

## Development

- `cd logstash && go test ./...`
