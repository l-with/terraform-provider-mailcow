package mailcow

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

import (
	"testing"
)

func TestAccDataSourceDkim(t *testing.T) {
	domain := "domain-with.440044.xyz"
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
		},
	})
}

func testAccDataSourceDkimSimple(domain string, length int) string {
	return fmt.Sprintf(`
resource "mailcow_domain" "domain" {
  domain = "%[1]s"
}

resource "mailcow_dkim" "dkim" {
  domain = mailcow_domain.domain.id
  length = %[2]d
}

data "mailcow_dkim" "demo" {
  domain = mailcow_dkim.dkim.id
}
`, domain, length)
}
