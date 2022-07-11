package main

import (
	"context"
	"terraform-provider-mailcow/mailcow"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

func main() {
	tfsdk.Serve(context.Background(), mailcow.New, tfsdk.ServeOpts{
		Name: "mailcow",
	})
}
