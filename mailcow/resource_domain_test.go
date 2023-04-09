package mailcow

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceDomain(t *testing.T) {
	percentS := "%s"
	domainFmt := fmt.Sprintf("with-domain-%s%s.domain-%s.xyz", randomLowerCaseString(4), percentS, randomLowerCaseString(4))
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDomainSimple(fmt.Sprintf(domainFmt, "1")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_domain.domain-simple", "id", fmt.Sprintf(domainFmt, "1")),
				),
			},
			{
				Config:   testAccResourceDomainSimple(fmt.Sprintf(domainFmt, "1")),
				PlanOnly: true,
			},
			{
				Config: testAccResourceDomainSimple(fmt.Sprintf(domainFmt, "2")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_domain.domain-simple", "id", fmt.Sprintf(domainFmt, "2")),
				),
			},
			{
				Config: testAccResourceDomainSimpleUpdate(fmt.Sprintf(domainFmt, "2")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_domain.domain-simple", "aliases", "1000"),
				),
			},
			{
				Config: testAccResourceDomain(fmt.Sprintf(domainFmt, "backup"), "true", "42000"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_domain.domain", "id", fmt.Sprintf(domainFmt, "backup")),
					resource.TestCheckResourceAttr("mailcow_domain.domain", "backupmx", "true"),
					resource.TestCheckResourceAttr("mailcow_domain.domain", "quota", "42000"),
				),
			},
			{
				Config: testAccResourceDomain(fmt.Sprintf(domainFmt, "backup"), "false", "84000"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_domain.domain", "id", fmt.Sprintf(domainFmt, "backup")),
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

func testAccResourceDomainSimple(domain string) string {
	return fmt.Sprintf(`
resource "mailcow_domain" "domain-simple" {
  domain   = "%[1]s"
}
`, domain)
}

func testAccResourceDomainSimpleUpdate(domain string) string {
	return fmt.Sprintf(`
resource "mailcow_domain" "domain-simple" {
  domain  = "%[1]s"
  aliases = 1000
}
`, domain)
}

func testAccResourceDomainCreateError() string {
	return `
resource "mailcow_domain" "domain-create-error" {
  domain = "%"
}
`
}

func testAccResourceDomain(domain string, backupmx string, quota string) string {
	return fmt.Sprintf(`
resource "mailcow_domain" "domain" {
  domain   = "%[1]s"
  backupmx = %[2]s
  quota = %[3]s
}
`, domain, backupmx, quota)
}
