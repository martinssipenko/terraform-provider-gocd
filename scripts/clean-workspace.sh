#!/usr/bin/env bash -x
go generate -x ./... && git diff --exit-code; code=$?; git checkout -- .; (exit $code)
