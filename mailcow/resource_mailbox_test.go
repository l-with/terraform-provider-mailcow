package mailcow

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceMailbox(t *testing.T) {
	domain := "domain-with4mailbox-test.440044.xyz"
	localPart := "mailbox-with"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceMailboxSimple(domain, localPart),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_mailbox.mailbox", "address", localPart+"@"+domain),
				),
			},
		},
	})
}

func testAccResourceMailboxSimple(domain string, localPart string) string {
	return fmt.Sprintf(`
resource "mailcow_domain" "domain" {
  domain = "%[1]s"
}

resource "mailcow_mailbox" "mailbox" {
  domain     = mailcow_domain.domain.domain 
  local_part = "%[2]s"
  password   = "password"
  quota      = 3072
}`, domain, localPart)
}
