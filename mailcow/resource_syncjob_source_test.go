package mailcow

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSyncjobSource(t *testing.T) {
	domain := fmt.Sprintf("with-syncjob-source-%s.domain-%s.xyz", randomLowerCaseString(4), randomLowerCaseString(4))
	name := fmt.Sprintf("with-syncjob-source-%s", randomLowerCaseString(4))
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSyncjobSource(domain, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_syncjob_source.source", "name", name),
					resource.TestCheckResourceAttr("mailcow_syncjob_source.source", "host", "imap.example.com"),
					resource.TestCheckResourceAttr("mailcow_syncjob_source.source", "port", "993"),
					resource.TestCheckResourceAttr("mailcow_syncjob_source.source", "enc", "SSL"),
					resource.TestCheckResourceAttr("mailcow_syncjob_source.source", "verify_cert", "true"),
					resource.TestCheckResourceAttr("mailcow_syncjob_source.source", "auth_type", "PLAIN"),
					resource.TestCheckResourceAttr("mailcow_syncjob_source.source", "scope", "user"),
					resource.TestCheckResourceAttr("mailcow_syncjob_source.source", "active", "true"),
					resource.TestCheckResourceAttrSet("mailcow_syncjob_source.source", "created_by"),
				),
			},
			{
				Config: testAccResourceSyncjobSourceUpdate(domain, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_syncjob_source.source", "name", name+"-update"),
					resource.TestCheckResourceAttr("mailcow_syncjob_source.source", "description", "updated"),
					resource.TestCheckResourceAttr("mailcow_syncjob_source.source", "host", "update-imap.example.com"),
					resource.TestCheckResourceAttr("mailcow_syncjob_source.source", "port", "143"),
					resource.TestCheckResourceAttr("mailcow_syncjob_source.source", "enc", "TLS"),
					resource.TestCheckResourceAttr("mailcow_syncjob_source.source", "verify_cert", "false"),
					resource.TestCheckResourceAttr("mailcow_syncjob_source.source", "auth_type", "LOGIN"),
					resource.TestCheckResourceAttr("mailcow_syncjob_source.source", "scope", "domain"),
					resource.TestCheckResourceAttr("mailcow_syncjob_source.source", "domains.#", "1"),
					resource.TestCheckTypeSetElemAttr("mailcow_syncjob_source.source", "domains.*", domain),
					resource.TestCheckResourceAttr("mailcow_syncjob_source.source", "active", "false"),
				),
			},
			{
				ResourceName:      "mailcow_syncjob_source.source",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      testAccResourceSyncjobSourceScope(domain, name, `scope = "domain"`),
				ExpectError: regexp.MustCompile(`needs at least one entry in domains`),
			},
			{
				Config:      testAccResourceSyncjobSourceScope(domain, name, `users = ["demo@example.com"]`+"\n  scope = \"all\""),
				ExpectError: regexp.MustCompile(`users can only be set with scope = "user"`),
			},
		},
	})
}

func testAccResourceSyncjobSource(domain string, name string) string {
	return fmt.Sprintf(`
resource "mailcow_domain" "domain" {
  domain = "%[1]s"
}

resource "mailcow_syncjob_source" "source" {
  name = "%[2]s"
  host = "imap.example.com"
}
`, domain, name)
}

func testAccResourceSyncjobSourceUpdate(domain string, name string) string {
	return fmt.Sprintf(`
resource "mailcow_domain" "domain" {
  domain = "%[1]s"
}

resource "mailcow_syncjob_source" "source" {
  name        = "%[2]s-update"
  description = "updated"
  host        = "update-imap.example.com"
  port        = 143
  enc         = "TLS"
  verify_cert = false
  auth_type   = "LOGIN"
  scope       = "domain"
  domains     = [mailcow_domain.domain.id]
  active      = false
}
`, domain, name)
}

func testAccResourceSyncjobSourceScope(domain string, name string, scope string) string {
	return fmt.Sprintf(`
resource "mailcow_domain" "domain" {
  domain = "%[1]s"
}

resource "mailcow_syncjob_source" "source" {
  name = "%[2]s"
  host = "imap.example.com"
  %[3]s
}
`, domain, name, scope)
}
