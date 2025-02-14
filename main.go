package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/killmeplz/terraform-provider-ansible-awx/provider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
