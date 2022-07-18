package mailcow

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceDkim(t *testing.T) {
	domain := "domain-with.440044.xyz"
	length := 2048
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDkimSimple(domain, length),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_dkim.dkim", "dkim_selector", "dkim"),
					resource.TestCheckResourceAttr("mailcow_dkim.dkim", "length", fmt.Sprint(length)),
					resource.TestCheckResourceAttr("mailcow_dkim.dkim", "id", domain),
				),
			},
		},
	})
}

func testAccResourceDkimSimple(domain string, length int) string {
	return fmt.Sprintf(`
resource "mailcow_domain" "domain" {
  domain = "%[1]s"
}

resource "mailcow_dkim" "dkim" {
  domain = "%[1]s"
  length = %[2]d
}
`, domain, length)
}
