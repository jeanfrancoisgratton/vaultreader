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
sudo touch /var/log/vaultreader.log
sudo chmod 664 /var/log/vaultreader.log
sudo chown 0:devops /var/log/vaultreader.log