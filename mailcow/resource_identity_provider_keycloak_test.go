package mailcow

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceIdentityProviderKeycloak(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIdentityProviderKeycloak("realm"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_identity_provider_keycloak.keycloak", "realm", "realm"),
				),
			},
			{
				Config: testAccResourceIdentityProviderKeycloak("realm_update"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_identity_provider_keycloak.keycloak", "realm", "realm_update"),
				),
			},
		},
	})
}

func testAccResourceIdentityProviderKeycloak(realm string) string {
	return fmt.Sprintf(`
resource "mailcow_identity_provider_keycloak" "keycloak" {
  realm         = "%[1]s"
  client_id     = "client_id"
  client_secret = "client_secret"
  redirect_url  = "redirect_url"
  server_url    = "server_url"
  version       = "version"
}
`, realm)
}
