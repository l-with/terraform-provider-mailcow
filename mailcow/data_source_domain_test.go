package mailcow

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDomain(t *testing.T) {
	domain := fmt.Sprintf("with-ds-domain-%s.domain-%s.xyz", randomLowerCaseString(4), randomLowerCaseString(4))
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
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
	return `
data "mailcow_domain" "error" {
  domain = "xyzzy"
}
`
}
