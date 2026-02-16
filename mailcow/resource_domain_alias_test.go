package mailcow

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceDomainAlias(t *testing.T) {
	percentS := "%s"
	aliasDomainFmt := fmt.Sprintf("alias-%s%s.domain-%s.xyz", randomLowerCaseString(4), percentS, randomLowerCaseString(4))
	targetDomainFmt := fmt.Sprintf("target-%s%s.domain-%s.xyz", randomLowerCaseString(4), percentS, randomLowerCaseString(4))

	aliasDomain := fmt.Sprintf(aliasDomainFmt, "1")
	targetDomain := fmt.Sprintf(targetDomainFmt, "1")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDomainAlias(aliasDomain, targetDomain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_domain_alias.test", "alias_domain", aliasDomain),
					resource.TestCheckResourceAttr("mailcow_domain_alias.test", "target_domain", targetDomain),
					resource.TestCheckResourceAttr("mailcow_domain_alias.test", "active", "true"),
				),
			},
			{
				Config: testAccResourceDomainAliasInactive(aliasDomain, targetDomain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_domain_alias.test", "active", "false"),
				),
			},
		},
	})
}

func testAccResourceDomainAlias(aliasDomain, targetDomain string) string {
	return fmt.Sprintf(`
resource "mailcow_domain" "target" {
  domain = "%s"
}

resource "mailcow_domain_alias" "test" {
  alias_domain  = "%s"
  target_domain = mailcow_domain.target.domain
}
`, targetDomain, aliasDomain)
}

func testAccResourceDomainAliasInactive(aliasDomain, targetDomain string) string {
	return fmt.Sprintf(`
resource "mailcow_domain" "target" {
  domain = "%s"
}

resource "mailcow_domain_alias" "test" {
  alias_domain  = "%s"
  target_domain = mailcow_domain.target.domain
  active        = false
}
`, targetDomain, aliasDomain)
}
