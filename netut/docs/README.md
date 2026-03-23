# netut Documentation

These notes split the repository-wide docs down to the `netut` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/netut`
- Run `cd netut && go test ./...` after changing network or TLS behavior.
- Keep transport-specific logic in the correct subpackage: `tcp`, `udp`, `https`, `hc`, `host`, `nettools`, or `ss`.
- Treat timeout and TLS changes as high-risk because several other modules depend on these helpers.

## API Surface
- `RetrieveIP`, `ParseIPs`, `GetLocalIP`, and `Ping` cover core IP utilities.
- `https.New`, `https.NewWithContext`, and `https.NewWithCustCA` set up HTTPS servers.
- `hc.NewClient`, client options, `MakeCall`, and `Reachable` build and execute HTTP requests.
- `ss.BasicAuth` wraps HTTP handlers with basic-auth and rotating-token checks.
- `nettools.CertificateInfo`, `CertificateChains`, and PEM helpers expose certificate utilities.

## Schema & Data Shapes
- No single external schema file is defined by the module.
- Main data shapes include `Host`, `hc.Client`, `hc.Auth`, `hc.ClientCert`, `HTTPClientSettings`, and `https.ServerCert`.

## Dataflow
1. Choose the protocol-specific helper that matches the network job.
2. Configure auth, TLS, timeouts, and request settings as needed.
3. Execute the request or server workflow and return typed network results to the caller.

## Operator Notes
- This is still a library module, but it often sits on live network paths.
- When debugging failures, inspect DNS resolution, TLS trust, timeout settings, and handler middleware ordering before changing code.
