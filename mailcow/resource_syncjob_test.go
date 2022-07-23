package mailcow

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSyncjob(t *testing.T) {
	domain := "domain-with4mailbox-test-440044.xyz"
	localPart := "localpart-with4mailbox-test"
	host1 := "example.com"
	user1 := "demo@example.com"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSyncjob(domain, localPart, host1, user1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_syncjob.syncjob", "username", localPart+"@"+domain),
					resource.TestCheckResourceAttr("mailcow_syncjob.syncjob", "delete2", "false"),
				),
			},
			{
				Config: testAccResourceSyncjobUpdate(domain, localPart, host1, user1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_syncjob.syncjob", "username", localPart+"@"+domain),
					resource.TestCheckResourceAttr("mailcow_syncjob.syncjob", "host1", "update-"+host1),
					resource.TestCheckResourceAttr("mailcow_syncjob.syncjob", "user1", "update-"+user1),
					resource.TestCheckResourceAttr("mailcow_syncjob.syncjob", "mins_interval", "42"),
					resource.TestCheckResourceAttr("mailcow_syncjob.syncjob", "maxbytespersecond", "42"),
					resource.TestCheckResourceAttr("mailcow_syncjob.syncjob", "delete2", "true"),
				),
			},
		},
	})
}

func testAccResourceSyncjob(domain string, localPart string, host1 string, user1 string) string {
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

resource "mailcow_syncjob" "syncjob" {
  username  = mailcow_mailbox.mailbox.address
  host1     = "%[3]s"
  user1     = "%[4]s"
  password1 = "secret-password"
}
`, domain, localPart, host1, user1)
}

func testAccResourceSyncjobUpdate(domain string, localPart string, host1 string, user1 string) string {
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

resource "mailcow_syncjob" "syncjob" {
  username          = mailcow_mailbox.mailbox.address
  host1             = "update-%[3]s"
  user1             = "update-%[4]s"
  password1         = "update-secret-password1"
  mins_interval     = 42
  maxbytespersecond = 42
  delete2           = true
}
`, domain, localPart, host1, user1)
}
