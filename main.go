package main

import (
	"github.com/hashicorp/terraform/plugin"

	"github.com/leonrodenburg/terraform-provider-alicloudext/alicloudext"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: alicloudext.Provider,
	})
}
