#!/bin/bash

# push.sh will only release when a new tag is 
# pushed to github and github actions runs the 
# push-release.yml action

LATEST_RELEASE=$(git describe --tags --abbrev=0)

echo "trying to find where we are in the directory bc i know nothing about github actions"
ls -al
pwd

# build both binaries
go build -o show-secrets ../
go build -o kubectl-show-secrets ../

echo "WE MADE IT THIS FAR"

# todo 
# update 