#!/bin/bash

# push.sh will only release when a new tag is 
# pushed to github and github actions runs the 
# push-release.yml action

LATEST_RELEASE=$(git tag --tags --abbrev=0)

# build both binaries
go build -o show-secrets ../
go build -o kubectl-show-secrets ../

echo "WE MADE IT THIS FAR"

# todo 
# update 