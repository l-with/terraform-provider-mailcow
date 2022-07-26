package mailcow

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
)

import (
	"testing"
)

func TestAccDataSourceDkim(t *testing.T) {
	domain := "domain-with4test-dkim.440044.xyz"
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
	return fmt.Sprintf(`
data "mailcow_dkim" "error" {
  domain = "xyzzy"
}
`)
}
