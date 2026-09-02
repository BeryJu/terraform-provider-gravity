package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceUser(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUserSimple(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gravity_user.test", "username", rName),
					resource.TestCheckResourceAttr("gravity_user.test", "permissions",
						`[{"methods":["GET"],"path":"/api/v1/dns/*"}]`),
				),
			},
		},
	})
}

// A user without permissions is returned as `null`, which has to normalise
// back to the `[]` the config asked for.
func TestAccResourceUserNoPermissions(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "gravity_user" "empty" {
  username = "%[1]s"
}
`, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gravity_user.empty", "permissions", "[]"),
				),
			},
		},
	})
}

func testAccResourceUserSimple(name string) string {
	return fmt.Sprintf(`
resource "gravity_user" "test" {
  username    = "%[1]s"
  password    = "acc-test-password"
  permissions = jsonencode([
    {
      path    = "/api/v1/dns/*"
      methods = ["GET"]
    }
  ])
}
`, name)
}
