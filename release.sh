#!/bin/sh

BUILD_VERSION=`git describe --tags`
COMMIT_HASH=`git rev-parse --short=8 HEAD 2>/dev/null`

gox -ldflags "-X main.Version=$BUILD_VERSION -X main.CommitHash=$COMMIT_HASH" -os="linux darwin windows openbsd" -arch="amd64" -output="dist/keyvault-get-secret_{{.OS}}_{{.Arch}}"