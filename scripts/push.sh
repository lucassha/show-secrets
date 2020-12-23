#!/bin/bash

# ./push.sh will only run when a new tag is 
# pushed to github and github actions runs the 
# push-release.yml action

set -eu

S3_TARBALL="showsecrets.tar.gz"
S3_BUCKET="lucassha-show-secrets-releases"
OS="darwin"
ARCH="amd64"

LATEST_RELEASE=$(git describe --tags --abbrev=0)

# build both binaries for mac
GOOS=${OS} GOARCH=${ARCH} go build -o show-secrets
GOOS=${OS} GOARCH=${ARCH} go build -o kubectl-show-secrets

tar -czvf ${S3_TARBALL}-${LATEST_RELEASE} ./kubectl-show-secrets ./show-secrets

aws s3 cp ${S3_TARBALL}-${LATEST_RELEASE} s3://${S3_BUCKET}/releases/${LATEST_RELEASE} --profile shannon