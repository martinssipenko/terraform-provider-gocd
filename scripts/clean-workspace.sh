#!/usr/bin/env bash
go generate -x ./... && git diff --exit-code; code=$?; git checkout -- .; (exit $code)
