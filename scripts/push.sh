#!/bin/bash

# ./push.sh will only run when a new tag is 
# pushed to github and github actions runs the 
# push-release.yml action

LATEST_RELEASE=$(git describe --tags --abbrev=0)

# build both binaries for mac
GOOS=darwin GOARCH=amd64 go build -o show-secrets
GOOS=darwin GOARCH=amd64 go build -o kubectl-show-secrets

tar -czvf showsecrets.tar.gz ./kubectl-show-secrets ./show-secrets

# todo 
# push up to S3