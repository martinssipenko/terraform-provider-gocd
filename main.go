package main

import (
	gocdprovider "github.com/beamly/terraform-provider-gocd/gocd"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: gocdprovider.Provider,
	})
}
