package mailcow

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceDomainAlias(t *testing.T) {
	aliasDomain := "alias-test.com"
	targetDomain := "test.com"

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
resource "mailcow_domain_alias" "test" {
  alias_domain  = "%s"
  target_domain = "%s"
}
`, aliasDomain, targetDomain)
}

func testAccResourceDomainAliasInactive(aliasDomain, targetDomain string) string {
	return fmt.Sprintf(`
resource "mailcow_domain_alias" "test" {
  alias_domain  = "%s"
  target_domain = "%s"
  active        = false
}
`, aliasDomain, targetDomain)
}
