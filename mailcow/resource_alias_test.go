package mailcow

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAlias(t *testing.T) {
	domain := fmt.Sprintf("with-alias-%s.domain-%s.xyz", randomLowerCaseString(4), randomLowerCaseString(4))
	localPart := fmt.Sprintf("with-alias-%s", randomLowerCaseString(4))
	percentS := "%s"
	aliasLocalPart := fmt.Sprintf("with-alias-%s-%s", randomLowerCaseString(4), percentS)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAliasSimple(fmt.Sprintf(aliasLocalPart, "1")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_alias.simple", "address", fmt.Sprintf(aliasLocalPart, "1")+"-simple@440044.xyz"),
				),
			},
			{
				Config: testAccResourceAliasSimple(fmt.Sprintf(aliasLocalPart, "1")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_alias.simple", "address", fmt.Sprintf(aliasLocalPart, "1")+"-simple@440044.xyz"),
					resource.TestCheckResourceAttr("mailcow_alias.simple", "sogo_visible", "false"),
				),
			},
			{
				Config: testAccResourceAliasSimpleUpdate(fmt.Sprintf(aliasLocalPart, "1")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_alias.simple", "sogo_visible", "true"),
				),
			},
			{
				Config: testAccResourceAliasSimpleUpdateToSpam(fmt.Sprintf(aliasLocalPart, "1")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_alias.simple", "goto", gotoSpamDestination),
				),
			},
			{
				Config: testAccResourceAliasSimpleUpdate(fmt.Sprintf(aliasLocalPart, "1")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_alias.simple", "goto", "demo@440044.xyz"),
				),
			},
			{
				Config: testAccResourceAliasSpecial(fmt.Sprintf(aliasLocalPart, "1"), "ham"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_alias.ham", "address", fmt.Sprintf(aliasLocalPart, "1")+"-ham@440044.xyz"),
					resource.TestCheckResourceAttr("mailcow_alias.ham", "goto", gotoHamDestination),
				),
			},
			{
				Config: testAccResourceAliasSpecial(fmt.Sprintf(aliasLocalPart, "1"), "null"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_alias.null", "address", fmt.Sprintf(aliasLocalPart, "1")+"-null@440044.xyz"),
					resource.TestCheckResourceAttr("mailcow_alias.null", "goto", gotoDiscardDestination),
				),
			},
			{
				Config: testAccResourceAliasSpecial(fmt.Sprintf(aliasLocalPart, "1"), "spam"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_alias.spam", "address", fmt.Sprintf(aliasLocalPart, "1")+"-spam@440044.xyz"),
					resource.TestCheckResourceAttr("mailcow_alias.spam", "goto", gotoSpamDestination),
				),
			},
			{
				Config: testAccResourceAlias(domain, localPart, fmt.Sprintf(aliasLocalPart, "2")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_alias.alias", "address", fmt.Sprintf(aliasLocalPart, "2")+"@"+domain),
				),
			},
			{
				Config:      testAccResourceAliasError("alias-xyzzy@xyzzy", "goto-xyzzy@xyzzy"),
				ExpectError: regexp.MustCompile("danger"),
			},
			{
				Config:      testAccResourceAliasUpdateError(domain, localPart, fmt.Sprintf(aliasLocalPart, "3")+"@"+domain),
				ExpectError: regexp.MustCompile("danger"),
			},
		},
	})
}

func TestAccResourceAlias_Validation(t *testing.T) {
	errorMessage := "cannot contain other addresses"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "mailcow_alias" "goto_ham_validation" {
						address   = "alias-xyzzy@xyzzy"
						goto      = "goto-xyzzy@xyzzy,ham@localhost"
					}
					`,
				ExpectError: regexp.MustCompile(errorMessage),
			},
			{
				Config: `
						resource "mailcow_alias" "goto_validation" {
							address   = "alias-xyzzy@xyzzy"
							goto      = "goto-xyzzy@xyzzy,null@localhost"
						}
						`,
				ExpectError: regexp.MustCompile(errorMessage),
			},
			{
				Config: `
						resource "mailcow_alias" "goto_validation" {
							address   = "alias-xyzzy@xyzzy"
							goto      = "goto-xyzzy@xyzzy,spam@localhost"
						}
						`,
				ExpectError: regexp.MustCompile(errorMessage),
			},
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

func testAccResourceAliasSpecial(aliasLocalPart, specialValue string) string {
	return fmt.Sprintf(`
resource "mailcow_alias" "%[1]s" {
  address = "%[2]s-%[1]s@440044.xyz"
  goto    = "%[1]s@localhost"
}
`, specialValue, aliasLocalPart)
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

func testAccResourceAliasSimpleUpdateToSpam(aliasLocalPart string) string {
	return fmt.Sprintf(`
resource "mailcow_alias" "simple" {
  address      = "%[1]s-demo@440044.xyz"
  goto         = "spam@localhost"
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
