#!/usr/bin/env sh

GOROOT=/opt/go
OUTPUT=/opt/bin
BINARY=vaultreader

# Get git branch's name, replace / with _
BRANCH=`git rev-parse --abbrev-ref HEAD`
BRANCH=$(echo "$BRANCH" | tr '/' '_')

if [ "$BRANCH" = "master" ] || [ "$BRANCH" = "main" ] || [ "$BRANCH" = "develop" ]; then
    FULLNAME="$BINARY"
else
    FULLNAME="$BINARY-$BRANCH"
fi

if [ "$#" -gt 0 ]; then
    OUTPUT=$1
fi

go build -trimpath -ldflags="-s -w -buildid=" -o $OUTPUT/$FULLNAME .
if getent group vaultreader > /dev/null 2>&1; then
    :  # group exists, nothing to do
else
    if getent group 3000 > /dev/null 2>&1; then
        groupadd vaultreader
    else
        groupadd -g 3000 vaultreader
    fi
fi

sudo touch /var/log/vaultreader.log
sudo chmod 664 /var/log/vaultreader.log
sudo chmod 2755 "$OUTPUT/$FULLNAME"
sudo chown 0:vaultreader /var/log/vaultreader.log "$OUTPUT/$FULLNAME"