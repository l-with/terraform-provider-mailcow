package mailcow

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

import (
	"testing"
)

func TestAccDataSourceStage(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDomainSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mailcow_domain.demo", "domain_name", "440044.xyz"),
				),
			},
		},
	})
}

const testAccDataSourceDomainSimple = `
data "mailcow_domain" "demo" {
  domain_name = "440044.xyz"
}
`
