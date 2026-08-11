package mailcow

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceRelayhost(t *testing.T) {
	hostname := fmt.Sprintf("relayhost-test-%s.xzy:%d", randomLowerCaseString(4), rand.Int())
	username := fmt.Sprintf("username-%s", randomLowerCaseString(4))
	password := fmt.Sprintf("password-%s", randomLowerCaseString(4))

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRelayhost(hostname, username, password),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_relayhost.relayhost", "hostname", hostname),
					resource.TestCheckResourceAttr("mailcow_relayhost.relayhost", "username", username),
					resource.TestCheckResourceAttr("mailcow_relayhost.relayhost", "password", password),
				),
			},
		},
	})
}

func testAccResourceRelayhost(domain string, username string, password string) string {
	return fmt.Sprintf(`
resource "mailcow_relayhost" "relayhost" {
	hostname = "%s"
	username = "%s"
	password = "%s"
}
`, domain, username, password)
}
