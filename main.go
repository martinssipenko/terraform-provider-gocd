package main

import (
	gocdprovider "github.com/drewsonne/terraform-provider-gocd/gocd"
	"github.com/hashicorp/terraform/plugin"
)

//func main() {
//	//plugin.Serve(gocdprovider.)
//	plugin.Serve(&plugin.ServeOpts{
//		ProviderFunc: gocdprovider.Provider(),
//	})
//}
func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: gocdprovider.Provider})
}
