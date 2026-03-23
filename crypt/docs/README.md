# crypt Documentation

These notes split the repository-wide docs down to the `crypt` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/crypt`
- Run `cd crypt && go test ./...` before and after algorithm changes.
- Keep algorithm-specific code in subpackages such as `encrypt`, `sign`, `verify`, `cert`, and `genppk`.
- Preserve compatibility expectations around key sizes, salts, nonces, and serialized hash formats.

## API Surface
- `encrypt/pbkdf2.New`, `encrypt/scrypt.Key`, and `encrypt/argon2id.Key` derive keys.
- `encrypt/aescrypt.New`, `encrypt/aescrypt_gcm.New`, and `encrypt/stream.Encrypt` handle encryption flows.
- `sign.NewRandomKey` and `sign.NewRandomNonce` support Poly1305 signing.
- `verify.NewRSAVerifier`, `verify.NewECDSAVerifier`, and `verify.NewDSAVerifier` expose signature verification.
- `cert.New` and `cert.LECerts` cover development and Let's Encrypt certificate workflows.

## Schema & Data Shapes
- There is no external config file schema in this module.
- Main runtime shapes include `sign.Key`, `sign.MACKey`, `encrypt/scrypt.Params`, `cert.Cert`, and the stream metadata structs.

## Dataflow
1. Derive or generate key material using the appropriate subpackage.
2. Encrypt, sign, or generate certificates with the matching API.
3. Verify, decrypt, or distribute the resulting material in the caller or downstream module.

## Operator Notes
- Library-only module; there is no daemon to run.
- When changing certificate generation or streaming encryption, validate with real files and not only unit tests.
