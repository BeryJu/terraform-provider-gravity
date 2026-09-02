package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAPIToken(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAPIToken(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gravity_api_token.test", "username", rName),
					resource.TestCheckResourceAttrSet("gravity_api_token.test", "key"),
				),
			},
		},
	})
}

func testAccResourceAPIToken(name string) string {
	return fmt.Sprintf(`
resource "gravity_user" "test" {
  username = "%[1]s"
  password = "acc-test-password"
}

resource "gravity_api_token" "test" {
  username = gravity_user.test.username
}
`, name)
}
