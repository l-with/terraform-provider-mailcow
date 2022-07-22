package mailcow

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

import (
	"testing"
)

func TestAccDataSourceMailbox(t *testing.T) {
	domain := "domain-with4mailbox-test-440044.xyz"
	localPart := "localpart-with4mailbox-data-test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMailbox(domain, localPart),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mailcow_mailbox.mailbox", "domain", domain),
					resource.TestCheckResourceAttr("data.mailcow_mailbox.mailbox", "local_part", localPart),
				),
			},
		},
	})
}

func testAccDataSourceMailbox(domain string, localPart string) string {
	return fmt.Sprintf(`
resource "mailcow_domain" "domain" {
  domain = "%[1]s"
}

resource "mailcow_mailbox" "mailbox" {
  local_part = "%[2]s"
  domain     = mailcow_domain.domain.id
  password   = "secret-password"
}

data "mailcow_mailbox" "mailbox" {
  address = mailcow_mailbox.mailbox.address
}
`, domain, localPart)
}
