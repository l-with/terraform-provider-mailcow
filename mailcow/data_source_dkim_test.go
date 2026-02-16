package mailcow

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDkim(t *testing.T) {
	domain := fmt.Sprintf("with-ds-dkim-%s.domain-%s.xyz", randomLowerCaseString(4), randomLowerCaseString(4))
	length := 2048
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDkimSimple(domain, length),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mailcow_dkim.demo", "id", domain),
					resource.TestCheckResourceAttr("data.mailcow_dkim.demo", "dkim_selector", "dkim"),
					resource.TestCheckResourceAttr("data.mailcow_dkim.demo", "length", fmt.Sprint(length)),
				),
			},
			{
				Config:      testAccDataSourceDkimError(),
				ExpectError: regexp.MustCompile("not found"),
			},
		},
	})
}

func testAccDataSourceDkimSimple(domain string, length int) string {
	return fmt.Sprintf(`
resource "mailcow_domain" "domain-dkim-ds" {
  domain = "%[1]s"
}

resource "mailcow_dkim" "dkim-ds" {
  domain = mailcow_domain.domain-dkim-ds.domain
  length = %[2]d
}

data "mailcow_dkim" "demo" {
  domain = mailcow_dkim.dkim-ds.domain
}
`, domain, length)
}

func testAccDataSourceDkimError() string {
	return `
data "mailcow_dkim" "error" {
  domain = "xyzzy"
}
`
}
