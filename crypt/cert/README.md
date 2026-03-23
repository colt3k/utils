# Alternative to using code when on Mac or Linux

This package contains two common certificate flows:

- `devcert.New(...)` for generating development certificates inside Go
- `LECerts(host, email)` for wiring Let’s Encrypt autocert into a server

## OpenSSL Shortcut

Generate a private key and self-signed certificate for localhost with this command:

```sh
openssl req -x509 -out localhost.crt -keyout localhost.key \
  -newkey rsa:2048 -nodes -sha256 \
  -subj '/CN=localhost' -extensions EXT -config <( \
  printf "[dn]\nCN=localhost\n[req]\ndistinguished_name = dn\n[EXT]\nsubjectAltName=DNS:localhost\nkeyUsage=digitalSignature\nextendedKeyUsage=serverAuth")
```

## Development

- `cd crypt && go test ./cert/...`
