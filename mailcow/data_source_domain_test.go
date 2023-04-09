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
	domain := fmt.Sprintf("with-ds-domain-%s.domain-%s.xyz", randomLowerCaseString(4), randomLowerCaseString(4))
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDomainSimple(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mailcow_domain.simple", "description", "demo domain"),
				),
			},
			{
				Config: testAccDataSourceDomain(domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mailcow_domain.demo", "description", "description"),
				),
			},
			{
				Config:      testAccDataSourceDomainError(),
				ExpectError: regexp.MustCompile("not found"),
			},
		},
	})
}

func testAccDataSourceDomainSimple() string {
	return fmt.Sprintf(`
data "mailcow_domain" "simple" {
  domain = "440044.xyz"
}
`)
}

func testAccDataSourceDomain(domain string) string {
	return fmt.Sprintf(`
resource "mailcow_domain" "domain" {
  domain = "%[1]s"
  description = "description"
}

data "mailcow_domain" "demo" {
  domain = mailcow_domain.domain.domain
}
`, domain)
}

func testAccDataSourceDomainError() string {
	return fmt.Sprintf(`
data "mailcow_domain" "error" {
  domain = "xyzzy"
}
`)
}
