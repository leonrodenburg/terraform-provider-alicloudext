package main

import (
	"github.com/hashicorp/terraform/plugin"

	"github.com/leonrodenburg/terraform-provider-alicloudssl/alicloudssl"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: alicloudssl.Provider,
	})
}
