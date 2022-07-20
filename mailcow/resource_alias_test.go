package mailcow

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAlias(t *testing.T) {
	domain := "domain-with4test-domain.440044.xyz"
	localPart := "localpart-with4alias-test"
	aliasLocalPart := "alias-localpart-with4alias-test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAliasSimple(domain, localPart, aliasLocalPart),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_alias.alias", "address", aliasLocalPart+"@"+domain),
				),
			},
			{
				Config: testAccResourceAliasSimple(domain, localPart, aliasLocalPart+"2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_alias.alias", "address", aliasLocalPart+"2@"+domain),
				),
			},
		},
	})
}

func testAccResourceAliasSimple(domain string, localPart string, aliasLocalPart string) string {
	return fmt.Sprintf(`
resource "mailcow_domain" "domain" {
  domain = "%[1]s"
}

resource "mailcow_mailbox" "mailbox" {
  local_part = "%[2]s"
  domain     = mailcow_domain.domain.id
  password   = "secret-password"
}

resource "mailcow_alias" "alias" {
  address = "%[3]s@${mailcow_domain.domain.domain}"
  goto    = mailcow_mailbox.mailbox.address
}
`, domain, localPart, aliasLocalPart)
}
