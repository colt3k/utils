# hash Documentation

These notes split the repository-wide docs down to the `hash` module. Start with [../README.md](../README.md) for the package overview.

## Developer Notes
- Module path: `github.com/colt3k/utils/hash`
- Run `cd hash && go test ./...` after changing formatting or algorithm behavior.
- Keep algorithm constructors consistent across `md`, `sha1`, `sha2`, `sha3`, `sha512`, and `blake2`.
- Review downstream modules such as `file`, `crypt`, and `updater` when serialized output changes.

## API Surface
- `Hasher` is the shared interface used by the algorithm wrappers.
- `hash.String` and `hash.File` hash raw strings or files with a standard `hash.Hash`.
- `*_hash.NewHash(opts...)` constructors configure algorithm-specific helpers.
- `h_mac.Hash`, `EncodeSecretAnsibleVault`, and `DecodeSecretAnsibleVault` cover HMAC and secret transport helpers.

## Schema & Data Shapes
- No external schema file is used.
- The main typed shapes are `hashenum.HashEnum`, algorithm-specific option functions, and `h_mac.Key` / `h_mac.Secret`.

## Dataflow
1. Choose a hash implementation or constructor.
2. Hash raw bytes, strings, files, or HMAC payloads.
3. Return the encoded or raw digest to the caller or downstream module.

## Operator Notes
- Library-only module; there is no runtime operator surface.
- Digest-format changes are high impact because many other modules depend on stable output.
