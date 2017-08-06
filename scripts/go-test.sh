#!/usr/bin/env bash -x -e

ROOT_DIR=$(pwd)/../
COVERAGE_PATH=${ROOT_DIR}/coverage.txt

echo "" > ${COVERAGE_PATH}

for d in $(go list ./... | grep -v vendor | grep -v gocd-response-links); do
    go test -v -race -coverprofile=profile.out -covermode=atomic $d
    if [ -f profile.out ]; then
        cat profile.out >> ${COVERAGE_PATH}
        rm profile.out
    fi
done