package mailcow

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOAuth2Client(t *testing.T) {
	redirectUri := fmt.Sprintf("https:/redirect%s.uri", randomLowerCaseString(4))
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOAuth2ClientSimple(redirectUri),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mailcow_oauth2_client.client", "redirect_uri", redirectUri),
					resource.TestCheckResourceAttr("mailcow_oauth2_client.client", "scope", "profile"),
					resource.TestMatchResourceAttr("mailcow_oauth2_client.client", "client_id", regexp.MustCompile("^[a-z0-9]{12}$")),
					resource.TestMatchResourceAttr("mailcow_oauth2_client.client", "client_secret", regexp.MustCompile("^[a-z0-9]{24}$")),
				),
			},
			{
				Config:      testAccResourceOAuth2ClientError(redirectUri),
				ExpectError: regexp.MustCompile("Error running apply"),
			},
		},
	})
}

func testAccResourceOAuth2ClientSimple(redirectUri string) string {
	return fmt.Sprintf(`
resource "mailcow_oauth2_client" "client" {
  redirect_uri = "%[1]s"
}
`, redirectUri)
}

func testAccResourceOAuth2ClientError(redirectUri string) string {
	return fmt.Sprintf(`
resource "mailcow_oauth2_client" "client1" {
  redirect_uri = "%[1]s"
}

resource "mailcow_oauth2_client" "client2" {
  depends_on = [ mailcow_oauth2_client.client1 ]
  redirect_uri = "%[1]s"
}
`, redirectUri)
}
