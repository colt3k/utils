# netut

Network Utilities

## Included Packages

- http client
- https server
- certificate tools
- tcp client/server
- udp client/server
- IP tools
- host lookup helpers
- session/basic auth middleware

## Common Entry Points

- `RetrieveIP`, `ParseIPs`, `GetLocalIP`, `Ping`
- `https.New`, `https.NewWithContext`, `https.NewWithCustCA`
- `hc.NewClient`, `processor.MakeCall`, `Reachable`
- `ss.BasicAuth`
- `nettools.CertificateInfo`, `CertificateChains`, `OutputPEMFile`

## Development

- `cd netut && go test ./...`
