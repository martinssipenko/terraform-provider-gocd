# terraform-provider-gocd 0.0.7

[![GoDoc](https://godoc.org/github.com/drewsonne/terraform-provider-gocd/gocd?status.svg)](https://godoc.org/github.com/drewsonne/terraform-provider-gocd/gocd)
[![Build Status](https://travis-ci.org/drewsonne/terraform-provider-gocd.svg?branch=master)](https://travis-ci.org/drewsonne/terraform-provider-gocd)
[![codecov](https://codecov.io/gh/drewsonne/terraform-provider-gocd/branch/master/graph/badge.svg)](https://codecov.io/gh/drewsonne/terraform-provider-gocd)
[![Go Report Card](https://goreportcard.com/badge/github.com/drewsonne/terraform-provider-gocd)](https://goreportcard.com/report/github.com/drewsonne/terraform-provider-gocd)

## Terraform provider

Terraform provider for GoCD Server

### Building the Provider

## Demo

You will need docker and terraform >= 0.10.0 installed for this demo to work.

Either build the provider with `go build` or download it from the gihub repository. If you download it, make sure the binary is in the current folder.

	$ go build

Spin up the test gocd server

    $ make provision-test-gocd

Then initialise and apply the configuration.

    $ terraform init
    $ terraform apply

