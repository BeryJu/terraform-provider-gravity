package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceDHCPLease(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDHCPLeaseSimple(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gravity_dhcp_scope.name", "name", rName),
					resource.TestCheckResourceAttr("gravity_dhcp_scope.name", "subnet_cidr", "10.10.10.0/24"),
					resource.TestCheckResourceAttr("gravity_dhcp_scope.name", "option.#", "2"),
					resource.TestCheckResourceAttr("gravity_dhcp_lease.record", "identifier", rName),
				),
			},
		},
	})
}

func testAccResourceDHCPLeaseSimple(name string) string {
	return fmt.Sprintf(`
resource "gravity_dhcp_scope" "name" {
  name          = "%[1]s"
  subnet_cidr = "10.10.10.0/24"
  ipam = {
    type = "internal"
    range_start = "10.10.10.100"
    range_end = "10.10.10.150"
  }
  option {
    tag_name = "router"
    value = "10.10.10.1"
  }
  option {
    tag_name = "name_server"
    // value = "10.10.10.1"
	value64 = [base64encode("10.10.10.1")]
  }
  dns {
	zone = "foo.bar."
	add_zone_in_hostname = true
  }
}

resource "gravity_dhcp_lease" "record" {
  scope = gravity_dhcp_scope.name.name
  hostname = "foo"
  address = "10.10.10.25"
  identifier = "%[1]s"
}
`, name)
}
