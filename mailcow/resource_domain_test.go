package mailcow

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceDomainOK(t *testing.T) {
	subdomainPrefix := "domain-with4domain-test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDomainSimple(subdomainPrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_domain.domain-simple", "id", subdomainPrefix+".440044.xyz"),
				),
			},
			{
				Config:   testAccResourceDomainSimple(subdomainPrefix),
				PlanOnly: true,
			},
			{
				Config: testAccResourceDomainSimple(subdomainPrefix + "2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_domain.domain-simple", "id", subdomainPrefix+"2.440044.xyz"),
				),
			},
			{
				Config: testAccResourceDomainSimpleUpdate(subdomainPrefix + "2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_domain.domain-simple", "aliases", "1000"),
				),
			},
			{
				Config: testAccResourceDomain(subdomainPrefix+"-backup", "true", "42000"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_domain.domain", "id", subdomainPrefix+"-backup.440044.xyz"),
					resource.TestCheckResourceAttr("mailcow_domain.domain", "backupmx", "true"),
					resource.TestCheckResourceAttr("mailcow_domain.domain", "quota", "42000"),
				),
			},
			{
				Config: testAccResourceDomain(subdomainPrefix+"-backup", "false", "84000"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_domain.domain", "id", subdomainPrefix+"-backup.440044.xyz"),
					resource.TestCheckResourceAttr("mailcow_domain.domain", "backupmx", "false"),
					resource.TestCheckResourceAttr("mailcow_domain.domain", "quota", "84000"),
				),
			},
			{
				Config:      testAccResourceDomainCreateError(),
				ExpectError: regexp.MustCompile("."),
			},
		},
	})
}

func testAccResourceDomainSimple(name string) string {
	return fmt.Sprintf(`
resource "mailcow_domain" "domain-simple" {
  domain   = "%[1]s.440044.xyz"
}
`, name)
}

func testAccResourceDomainSimpleUpdate(name string) string {
	return fmt.Sprintf(`
resource "mailcow_domain" "domain-simple" {
  domain  = "%[1]s.440044.xyz"
  aliases = 1000
}
`, name)
}

func testAccResourceDomainCreateError() string {
	return `
resource "mailcow_domain" "domain-create-error" {
  domain = "%"
}
`
}

func testAccResourceDomain(prefix string, backupmx string, quota string) string {
	return fmt.Sprintf(`
resource "mailcow_domain" "domain" {
  domain   = "%[1]s.440044.xyz"
  backupmx = %[2]s
  quota = %[3]s
}
`, prefix, backupmx, quota)
}
