package mailcow

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAlias(t *testing.T) {
	aliasPrefix := "alias-with"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAliasSimple(aliasPrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_alias.alias", "address", aliasPrefix+"-demo@440044.xyz"),
				),
			},
			{
				Config: testAccResourceAliasSimple(aliasPrefix + "2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_alias.alias", "address", aliasPrefix+"2-demo@440044.xyz"),
				),
			},
		},
	})
}

func testAccResourceAliasSimple(name string) string {
	return fmt.Sprintf(`
resource "mailcow_alias" "alias" {
  address = "%[1]s-demo@440044.xyz"
  goto = "demo@440044.xyz"
}
`, name)
}
