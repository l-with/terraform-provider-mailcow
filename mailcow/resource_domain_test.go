package mailcow

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceDomain(t *testing.T) {
	subdomainPrefix := "domain-with4domain-test"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDomainSimple(subdomainPrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_domain.domain", "id", subdomainPrefix+".440044.xyz"),
				),
			},
			{
				Config: testAccResourceDomainSimple(subdomainPrefix + "2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_domain.domain", "id", subdomainPrefix+"2.440044.xyz"),
				),
			},
		},
	})
}

func testAccResourceDomainSimple(name string) string {
	return fmt.Sprintf(`
resource "mailcow_domain" "domain" {
  domain = "%[1]s.440044.xyz"
}
`, name)
}
