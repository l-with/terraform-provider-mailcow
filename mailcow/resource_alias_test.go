package mailcow

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAlias(t *testing.T) {
	//	domain := "domain-with4test-domain.440044.xyz"
	//	localPart := "localpart-with4alias-test"
	aliasLocalPart := "alias-localpart-with4alias-test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAliasSimple(aliasLocalPart),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_alias.simple", "address", aliasLocalPart+"-simple@440044.xyz"),
				),
			},
			/*
				{
					Config: testAccResourceAliasSimple(aliasLocalPart + "2"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("mailcow_alias.simple", "address", aliasLocalPart+"2-demo@440044.xyz"),
						resource.TestCheckResourceAttr("mailcow_alias.simple", "sogo_visible", "false"),
					),
				},
				{
					Config: testAccResourceAliasSimpleUpdate(aliasLocalPart + "2"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("mailcow_alias.simple", "sogo_visible", "true"),
					),
				},
				{
					Config: testAccResourceAlias(domain, localPart, aliasLocalPart),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("mailcow_alias.alias", "address", aliasLocalPart+"@"+domain),
					),
				},
				{
					Config: testAccResourceAlias(domain, localPart, aliasLocalPart+"2"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("mailcow_alias.alias", "address", aliasLocalPart+"2@"+domain),
					),
				},
				{
					Config:      testAccResourceAliasError("alias-xyzzy@xyzzy", "goto-xyzzy@xyzzy"),
					ExpectError: regexp.MustCompile("danger"),
				},
				{
					Config:      testAccResourceAliasUpdateError(domain, localPart, aliasLocalPart+"@"+domain),
					ExpectError: regexp.MustCompile("danger"),
				},

			*/
		},
	})
}

func testAccResourceAliasSimple(aliasLocalPart string) string {
	return fmt.Sprintf(`
resource "mailcow_alias" "simple" {
  address = "%[1]s-simple@440044.xyz"
  goto    = "demo@440044.xyz"
}
`, aliasLocalPart)
}

func testAccResourceAliasSimpleUpdate(aliasLocalPart string) string {
	return fmt.Sprintf(`
resource "mailcow_alias" "simple" {
  address      = "%[1]s-demo@440044.xyz"
  goto         = "demo@440044.xyz"
  sogo_visible = true
}
`, aliasLocalPart)
}

func testAccResourceAlias(domain string, localPart string, aliasLocalPart string) string {
	return fmt.Sprintf(`
resource "mailcow_domain" "domain" {
  domain = "%[1]s"
}

resource "mailcow_mailbox" "mailbox" {
  local_part = "%[2]s"
  domain     = mailcow_domain.domain.id
  password   = "secret-password"
  full_name  = "%[2]s"
}

resource "mailcow_alias" "alias" {
  address = "%[3]s@${mailcow_domain.domain.domain}"
  goto    = mailcow_mailbox.mailbox.address
}
`, domain, localPart, aliasLocalPart)
}

func testAccResourceAliasUpdateError(domain string, localPart string, aliasLocalPart string) string {
	return fmt.Sprintf(`
resource "mailcow_domain" "domain" {
  domain = "%[1]s"
}

resource "mailcow_mailbox" "mailbox" {
  local_part = "%[2]s"
  domain     = mailcow_domain.domain.id
  password   = "secret-password"
  full_name  = "%[2]s"
}

resource "mailcow_alias" "alias" {
  address = "%[3]s@$xyzzy"
  goto    = mailcow_mailbox.mailbox.address
}
`, domain, localPart, aliasLocalPart)
}

func testAccResourceAliasError(address string, gotoAddress string) string {
	return fmt.Sprintf(`
resource "mailcow_alias" "error" {
  address = "%[1]s"
  goto    = "%[2]s"
}
`, address, gotoAddress)
}
