# vaultreader

Read-only Vault client for KV v2 secrets. Almost always invoked from shell scripts and container entrypoints as `value=$(vaultreader -q ... -f field)`, which dictates everything below.

## stdout is the payload

The secret value, and nothing else, goes to stdout. **Every diagnostic goes to stderr**, via customError's `Print()` / `Die()` — never `fmt.Println(cerr.Error())`.

This is not a style preference. When Vault was sealed, the error text used to be captured *as the secret* by the calling `$(...)`; Jenkins received a skull-and-bones string as its keystore password and failed to start. The value was non-empty, so the caller's `[ -z ]` guard passed and the failure surfaced far from its cause.

`-q` suppresses the *payload* rendering, never the diagnostics.

This tool is the reference implementation for the customError v3.1.0 migration — see `docs/MIGRATION-v3.1.0.md` in the customError repo.

## Error codes

`types.ErrorMessages` maps each condition to a code and a stable `ERR_*` string. Set a real `Code`/`PosixErrorCode` on the error so calling scripts can tell a sealed Vault (retry, it is a boot-time race) from a bad token (give up). Do not add a condition to the iota list without also setting it somewhere — `ErrVaultSealed` and `ErrVaultUnavailable` sat declared-but-unused for a long time, which is why the sealed case was indistinguishable from any other failure.

## Callers to keep in mind

`docker_artifacts/jenkins/files/entrypoint.sh` is the most exposed consumer. Changing the output format, the exit codes, or the meaning of `-q` will affect container startup there.
