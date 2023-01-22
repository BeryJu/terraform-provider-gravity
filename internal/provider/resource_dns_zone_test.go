package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDNSZoneResource(t *testing.T) {
	zoneName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSZoneResourceConfig(zoneName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("gravity_dns_zone.test", "zone", zoneName),
					resource.TestCheckResourceAttr("gravity_dns_zone.test", "id", zoneName),
				),
			},
		},
	})
}

func testAccDNSZoneResourceConfig(zoneName string) string {
	return fmt.Sprintf(`
resource "gravity_dns_zone" "test" {
  zone = %[1]q

  handlers = [
	{
		type = "etcd",
		config = {},
	}
  ]
}
`, zoneName)
}
