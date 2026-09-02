package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceDHCPScope(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDHCPScopeSimple(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gravity_dhcp_scope.name", "name", rName),
					resource.TestCheckResourceAttr("gravity_dhcp_scope.name", "subnet_cidr", "10.10.10.0/24"),
					resource.TestCheckResourceAttr("gravity_dhcp_scope.name", "option.#", "2"),
					resource.TestCheckResourceAttr("gravity_dhcp_scope.name", "hook", "// noop"),
					resource.TestCheckResourceAttrSet("gravity_dhcp_scope.name", "statistics.0.usable"),
				),
			},
		},
	})
}

// Options that differ only in fields the old set hash could not read used to
// collapse into a single entry, leaving a permanent diff.
func TestAccResourceDHCPScopeMultipleOptions(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDHCPScopeMultipleOptions(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gravity_dhcp_scope.name", "option.#", "3"),
				),
			},
			{
				ResourceName:      "gravity_dhcp_scope.name",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceDHCPScopeMultipleOptions(name string) string {
	return fmt.Sprintf(`
resource "gravity_dhcp_scope" "name" {
  name        = "%[1]s"
  subnet_cidr = "10.11.12.0/24"
  ipam = {
    type        = "internal"
    range_start = "10.11.12.100"
    range_end   = "10.11.12.150"
  }
  option {
    tag_name = "router"
    value    = "10.11.12.1"
  }
  option {
    tag_name = "name_server"
    value    = "10.11.12.2"
  }
  option {
    tag_name = "domain_name"
    value    = "example.com"
  }
}
`, name)
}

func testAccResourceDHCPScopeSimple(name string) string {
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
  hook = "// noop"
}
`, name)
}
