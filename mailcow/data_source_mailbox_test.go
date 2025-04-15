package mailcow

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMailbox(t *testing.T) {
	domain := fmt.Sprintf("with-ds-mailbox-%s.domain-%s.xyz", randomLowerCaseString(4), randomLowerCaseString(4))
	localPart := fmt.Sprintf("with-ds-mailbox-%s", randomLowerCaseString(4))
	fullName := "full name"
	authSource := "mailcow"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMailbox(domain, localPart, fullName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mailcow_mailbox.mailbox", "domain", domain),
					resource.TestCheckResourceAttr("data.mailcow_mailbox.mailbox", "local_part", localPart),
					resource.TestCheckResourceAttr("data.mailcow_mailbox.mailbox", "full_name", fullName),
					resource.TestCheckResourceAttr("data.mailcow_mailbox.mailbox", "authsource", authSource),
				),
			},
			{
				Config:      testAccDataSourceMailboxError(),
				ExpectError: regexp.MustCompile("not found"),
			},
		},
	})
}

func testAccDataSourceMailbox(domain string, localPart string, fullName string) string {
	return fmt.Sprintf(`
resource "mailcow_domain" "domain" {
  domain = "%[1]s"
}

resource "mailcow_mailbox" "mailbox" {
  local_part = "%[2]s"
  domain     = mailcow_domain.domain.id
  password   = "secret-password"
  full_name  = "%[3]s"
}

data "mailcow_mailbox" "mailbox" {
  address = mailcow_mailbox.mailbox.address
}
`, domain, localPart, fullName)
}

func testAccDataSourceMailboxError() string {
	return fmt.Sprintf(`
data "mailcow_mailbox" "mailbox" {
  address = "xyzzy"
}
`)
}
