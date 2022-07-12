package main

import (
	"github.com/l-with/terraform-provider-mailcow/mailcow"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return mailcow.Provider()
		},
	})
}
