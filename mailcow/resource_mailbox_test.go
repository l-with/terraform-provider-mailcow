package mailcow

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceMailbox(t *testing.T) {
	domain := "domain-with4mailbox-test-440044.xyz"
	localPart := "localpart-with4mailbox-test"
	fullName := "new full name"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceMailbox(domain, localPart),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_mailbox.mailbox", "address", localPart+"@"+domain),
					resource.TestCheckResourceAttr("mailcow_mailbox.mailbox", "tls_enforce_out", "true"),
				),
			},
			{
				Config: testAccResourceMailboxUpdate(domain, localPart, fullName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_mailbox.mailbox", "full_name", fullName),
					resource.TestCheckResourceAttr("mailcow_mailbox.mailbox", "tls_enforce_out", "true"),
					resource.TestCheckResourceAttr("mailcow_mailbox.mailbox", "quota", "42"),
				),
			},
			{
				Config:      testAccResourceMailboxCreateError("xyzzy"),
				ExpectError: regexp.MustCompile("danger"),
			},
		},
	})
}

func testAccResourceMailbox(domain string, localPart string) string {
	return fmt.Sprintf(`
resource "mailcow_domain" "domain" {
  domain = "%[1]s"
}

resource "mailcow_mailbox" "mailbox" {
  local_part      = "%[2]s"
  domain          = mailcow_domain.domain.id
  password        = "secret-password"
  full_name       = "initial full name"
  tls_enforce_out = true
}
`, domain, localPart)
}

func testAccResourceMailboxUpdate(domain string, localPart string, fullName string) string {
	return fmt.Sprintf(`
resource "mailcow_domain" "domain" {
  domain = "%[1]s"
}

resource "mailcow_mailbox" "mailbox" {
  local_part      = "%[2]s"
  domain          = mailcow_domain.domain.id
  password        = "secret-password"
  full_name       = "%[3]s"
  tls_enforce_out = true
  quota           = 42
}
`, domain, localPart, fullName)
}

func testAccResourceMailboxCreateError(domain string) string {
	return `
resource "mailcow_mailbox" "mailbox-create" {
  local_part = "localpart"
  domain     = "%[1]s"
  password   = "secret-password"
  full_name  = "full name"
}
`
}
