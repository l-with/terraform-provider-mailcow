package mailcow

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceMailbox(t *testing.T) {
	domain := fmt.Sprintf("with-mailbox-%s.domain-%s.xyz", randomLowerCaseString(4), randomLowerCaseString(4))
	localPart := fmt.Sprintf("with-mailbox-%s", randomLowerCaseString(4))
	fullName := "new full name"
	quota := 4096
	domainMaxQuota := 5120
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceMailbox(domain, localPart),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_mailbox.mailbox", "address", localPart+"@"+domain),
					resource.TestCheckResourceAttr("mailcow_mailbox.mailbox", "tls_enforce_out", "true"),
					resource.TestCheckResourceAttr("mailcow_mailbox.mailbox", "authsource", mailcowAuthsourceInternal),
				),
			},
			{
				Config: testAccResourceMailboxUpdate(domain, domainMaxQuota, localPart, fullName, quota),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_mailbox.mailbox", "full_name", fullName),
					resource.TestCheckResourceAttr("mailcow_mailbox.mailbox", "tls_enforce_out", "true"),
					resource.TestCheckResourceAttr("mailcow_mailbox.mailbox", "quota", strconv.Itoa(quota)),
				),
			},
			{
				Config: testAccResourceMailboxUpdate(domain, domainMaxQuota, localPart, fullName, quota+1024),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_mailbox.mailbox", "full_name", fullName),
					resource.TestCheckResourceAttr("mailcow_mailbox.mailbox", "tls_enforce_out", "true"),
					resource.TestCheckResourceAttr("mailcow_mailbox.mailbox", "quota", strconv.Itoa(quota+1024)),
				),
			},
			{
				Config: testAccResourceMailboxChangeAuthsource(domain, domainMaxQuota, localPart, fullName, quota+1024),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_mailbox.mailbox", "authsource", mailcowAuthsourceKeycloak),
				),
			},
			{
				Config: testAccResourceMailboxAuthsource(domain, localPart),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_mailbox.mailbox_sso", "address", localPart+"-sso@"+domain),
					resource.TestCheckResourceAttr("mailcow_mailbox.mailbox_sso", "tls_enforce_out", "true"),
					resource.TestCheckResourceAttr("mailcow_mailbox.mailbox_sso", "authsource", mailcowAuthsourceKeycloak),
				),
			},
			{
				Config:      testAccResourceMailboxUpdate(domain, domainMaxQuota, localPart, fullName, domainMaxQuota+1024),
				ExpectError: regexp.MustCompile("danger"),
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

func testAccResourceMailboxUpdate(domain string, domainMaxquota int, localPart string, fullName string, quota int) string {
	return fmt.Sprintf(`
resource "mailcow_domain" "domain" {
  domain = "%[1]s"
  maxquota = %[2]d
}

resource "mailcow_mailbox" "mailbox" {
  local_part      = "%[3]s"
  domain          = mailcow_domain.domain.id
  password        = "secret-password"
  full_name       = "%[4]s"
  tls_enforce_out = true
  quota           = %[5]d
}
`, domain, domainMaxquota, localPart, fullName, quota)
}

func testAccResourceMailboxChangeAuthsource(domain string, domainMaxquota int, localPart string, fullName string, quota int) string {
	return fmt.Sprintf(`
resource "mailcow_domain" "domain" {
  domain = "%[1]s"
  maxquota = %[2]d
}

resource "mailcow_mailbox" "mailbox" {
  local_part      = "%[3]s"
  domain          = mailcow_domain.domain.id
  password        = "secret-password"
	force_pw_update = false
	authsource      = "keycloak"
  full_name       = "%[4]s"
  tls_enforce_out = true
  quota           = %[5]d
}
`, domain, domainMaxquota, localPart, fullName, quota)
}

func testAccResourceMailboxAuthsource(domain string, localPart string) string {
	return fmt.Sprintf(`
resource "mailcow_domain" "domain" {
  domain = "%[1]s"
}

resource "mailcow_mailbox" "mailbox_sso" {
  local_part      = "%[2]s-sso"
  domain          = mailcow_domain.domain.id
  password        = "secret-password"
	force_pw_update = false
	authsource      = "keycloak"
  full_name       = "initial full name"
  tls_enforce_out = true
}
`, domain, localPart)
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
