package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceDiscoverySubnet(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDiscoverySubnet(rName, "10.20.30.0/24"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gravity_discovery_subnet.test", "name", rName),
					resource.TestCheckResourceAttr("gravity_discovery_subnet.test", "subnet_cidr", "10.20.30.0/24"),
					resource.TestCheckResourceAttr("gravity_discovery_subnet.test", "dns_resolver", "10.20.30.1:53"),
					resource.TestCheckResourceAttr("gravity_discovery_subnet.test", "discovery_ttl", "3600"),
				),
			},
			{
				Config: testAccResourceDiscoverySubnet(rName, "10.20.31.0/24"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gravity_discovery_subnet.test", "subnet_cidr", "10.20.31.0/24"),
				),
			},
			{
				ResourceName:      "gravity_discovery_subnet.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceDiscoverySubnet(name string, cidr string) string {
	return fmt.Sprintf(`
resource "gravity_discovery_subnet" "test" {
  name          = "%[1]s"
  subnet_cidr   = "%[2]s"
  dns_resolver  = "10.20.30.1:53"
  discovery_ttl = 3600
}
`, name, cidr)
}
