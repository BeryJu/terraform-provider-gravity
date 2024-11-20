package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceDNSZone(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDNSZoneSimple(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gravity_dns_zone.name", "name", fmt.Sprintf("%s.", rName)),
					resource.TestCheckResourceAttr("gravity_dns_zone.name", "authoritative", "true"),
					resource.TestCheckResourceAttr("gravity_dns_zone.name", "handler_configs", "[{\"to\":[\"1.1.1.1\"],\"type\":\"etcd\"}]"),
				),
			},
		},
	})
}

func testAccResourceDNSZoneSimple(name string) string {
	return fmt.Sprintf(`
resource "gravity_dns_zone" "name" {
  name            = "%[1]s."
  authoritative   = true
  handler_configs = jsonencode([
    {
      type = "etcd",
	  to = ["1.1.1.1"],
    }
  ])
}
`, name)
}
