package mailcow

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

import (
	"testing"
)

func TestAccDataSourceDkim(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDkimSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mailcow_dkim.demo", "id", "440044.xyz"),
					resource.TestCheckResourceAttr("data.mailcow_dkim.demo", "dkim_selector", "dkim"),
					resource.TestCheckResourceAttr("data.mailcow_dkim.demo", "length", "2048"),
				),
			},
		},
	})
}

const testAccDataSourceDkimSimple = `
data "mailcow_dkim" "demo" {
  domain = "440044.xyz"
}
`
