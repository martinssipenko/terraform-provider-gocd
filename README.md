# terraform-provider-gocd 0.1.8

[![GoDoc](https://godoc.org/github.com/drewsonne/terraform-provider-gocd/gocd?status.svg)](https://godoc.org/github.com/drewsonne/terraform-provider-gocd/gocd)
[![Build Status](https://travis-ci.org/drewsonne/terraform-provider-gocd.svg?branch=master)](https://travis-ci.org/drewsonne/terraform-provider-gocd)
[![codecov](https://codecov.io/gh/drewsonne/terraform-provider-gocd/branch/master/graph/badge.svg)](https://codecov.io/gh/drewsonne/terraform-provider-gocd)
[![Go Report Card](https://goreportcard.com/badge/github.com/drewsonne/terraform-provider-gocd)](https://goreportcard.com/report/github.com/drewsonne/terraform-provider-gocd)

## Terraform provider
Terraform provider for GoCD Server

## Installation

    $ brew tap drewsonne/tap
    $ brew install terraform-provider-gocd
    $ tf-install-provider gocd
    
__NOTE__: `terraform` does not currently provide a way to easily install 3rd party providers. Until this is implemented,
the `tf-install-provider` utility can be used to copy the provider binary to the correct location.

### Building the Provider

## Demo

You will need docker and terraform >= 0.10.0 installed for this demo to work.

Either build the provider with `go build` or download it from the gihub repository. If you download it, make sure the binary is in the current folder.

	$ go build

Spin up the test gocd server, with endpoint at http://127.0.0.1:8153/go/

    $ make provision-test-gocd && sh ./scripts/wait-for-test-server.sh

Then initialise and apply the configuration.

    $ terraform init
    $ terraform apply

When you're finished, run:

    $ make teardown-test-gocd