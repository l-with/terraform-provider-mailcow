package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/l-with/terraform-provider-mailcow/internal"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: mailcow.Provider})
}
