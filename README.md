# vaultreader
___
A lightweight Hashicorp Vault secret reader.


# Introduction

This tool reads secrets off a Hashicorp Vault secret engine, specifically, a KV engine v2+. The engine thus supports versionned KV secrets.
___
# Pre-requisites

The following environment variables are not a hard requirement in the sense that they can be over-ridden with flags, in parentheseses.

## VAULT_ADDR (`-a $address`)
The actual server's address

## VAULT_AUTH_TOKEN (`-t $token`)
The authorization token to be authenticate against the Vault service.

## Optional: Quiet (`-q`)
This will suppress the output; useful when the tool is used in a CI/CD toolchain

## Optional: TEXT or JSON output (`-o {text|json}`)
If `-q` is invoked this option is ignored. If `-q` is not invoked, the output defaults to `-o text`

___
# How to use:
Simple: `vaultreader KV_ENGINE [-a vaultserver] [-t auth_token] [-v secret version number] [-f field name] secret_path`

- If `-v` is omitted, the latest secret's version will be read
- if `-f` is omitted, all fields in the secret will be fetched
___
# Building the software:
Whichever method you choose, go must be installed on the build machine (... well of course !). Check the file `go.version`
## From source:
Again, simple:
1. `cd $REPODIR/src`
2. (optional): run `./updateBuildDeps.sh`, to update all packages before compiling
3. `./build.sh`
By default it will compile the tool as `/opt/bin/vaultreader`, unless you over-ride this with
`./build.sh $OUTPUT`, where the binary would be compiled as $OUTPUT.

## Binary packages
Assuming you meet each distros' build framework requirements

### Alpine APK:
1. Run `abuild -r` from the `__alpine/` directory
2. The resulting .apk package will be located in `/data/packages`

### Debian DEB:
1. Go into the `__debian` directory and run the first two numbered shell scripts found there.
2. Once you've copied the resulting .deb package, you can run the last numbered script

### RedHat RPM:
The `tito` build tools must be installed on the build machine on top of the usual RPM build tools

1. If it is the first time you've built the rpm: `tito init`
2. `tito tag --keep-version`
3. `git push --follow-tags`
4. `tito build --rpm`
The resulting .rpm package will be found under `/tmp/tito`


