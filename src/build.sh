#!/usr/bin/env sh

set -e

BRANCH=`git rev-parse --abbrev-ref HEAD`
BRANCH=$(echo "$BRANCH" | tr '/' '_')
BINARY=vaultreader
OUTPUT=/opt/bin
COMPLETION=false

# Parse arguments
while [ "$#" -gt 0 ]; do
    case "$1" in
        -b|--binary)
            shift
            BINARY="$1"
            ;;
        *)
            OUTPUT="$1"
            ;;
    esac
    shift
done

if [ "$BRANCH" = "master" ] || [ "$BRANCH" = "main" ] || [ "$BRANCH" = "develop" ]; then
    FULLNAME="$BINARY"
else
    FULLNAME="$BINARY-$BRANCH"
fi

# set -e is enabled above, so both of these fast-fail the build; go test
# exits 0 for packages with no test files and only fails on an actual test
# failure.
go vet ./...
go test ./...

SCRIPT_DIR="$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)"
VRVERSION="$(sed -n 's/.*"versionnumber": *"\([^"]*\)".*/\1/p' "$SCRIPT_DIR/../vaultreader.json")"
BUILDDATE="$(date +%Y.%m.%d)"

echo "Building ${OUTPUT}/${FULLNAME}"
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -buildid= -X vaultreader/cmd.buildVersion=$VRVERSION -X vaultreader/cmd.buildDate=$BUILDDATE" -o ${OUTPUT}/${FULLNAME} .


# Enable tab completion

