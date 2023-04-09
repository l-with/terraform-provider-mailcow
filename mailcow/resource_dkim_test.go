package mailcow

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceDkim(t *testing.T) {
	domain := fmt.Sprintf("with-dkim-%s.dkim-%s.xyz", randomLowerCaseString(4), randomLowerCaseString(4))
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
					resource.TestMatchResourceAttr("mailcow_dkim.dkim", "dkim_txt", regexp.MustCompile("v=DKIM")),
				),
			},
			{
				Config:      testAccResourceDkimCreateError("xyzzy", length),
				ExpectError: regexp.MustCompile("danger"),
			},
		},
	})
}

func testAccResourceDkimSimple(domain string, length int) string {
	return fmt.Sprintf(`
resource "mailcow_domain" "domain-dkim" {
  domain = "%[1]s"
}

resource "mailcow_dkim" "dkim" {
  domain = mailcow_domain.domain-dkim.id
  length = %[2]d
}
`, domain, length)
}

func testAccResourceDkimCreateError(domain string, length int) string {
	return fmt.Sprintf(`
resource "mailcow_dkim" "error" {
  domain = "%[1]s"
  length = %[2]d
}
`, domain, length)
}
