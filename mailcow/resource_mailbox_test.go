package mailcow

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceMailbox(t *testing.T) {
	domain := fmt.Sprintf("with-mailbox-%s.domain-%s.xyz", randomLowerCaseString(4), randomLowerCaseString(4))
	localPart := fmt.Sprintf("with-mailbox-%s", randomLowerCaseString(4))
	fullName := "new full name"
	quota := 42
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
				Config: testAccResourceMailboxUpdate(domain, localPart, fullName, quota),
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

func testAccResourceMailboxUpdate(domain string, localPart string, fullName string, quota int) string {
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
  quota           = %[4]d
}
`, domain, localPart, fullName, quota)
}

func testAccResourceMailboxCreateError(domain string) string {
	return fmt.Sprintf(`
resource "mailcow_mailbox" "mailbox-create" {
  local_part = "localpart"
  domain     = "%[1]s"
  password   = "secret-password"
  full_name  = "full name"
}
`, domain)
}
