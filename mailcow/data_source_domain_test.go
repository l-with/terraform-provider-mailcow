package mailcow

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
)

import (
	"testing"
)

func TestAccDataSourceDomain(t *testing.T) {
	domain := "domain-with4test-domain.440044.xyz"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDomainSimple(domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mailcow_domain.demo", "domain_name", domain),
				),
			},
			{
				Config:      testAccDataSourceDomainError(),
				ExpectError: regexp.MustCompile("not found"),
			},
		},
	})
}

func testAccDataSourceDomainSimple(domain string) string {
	return fmt.Sprintf(`
resource "mailcow_domain" "domain" {
  domain = "%[1]s"
}

data "mailcow_domain" "demo" {
  domain_name = mailcow_domain.domain.domain
}
`, domain)
}

func testAccDataSourceDomainError() string {
	return fmt.Sprintf(`
data "mailcow_domain" "error" {
  domain_name = "xyzzy"
}
`)
}
