package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceDNSRecord(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDNSRecordSimple(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gravity_dns_zone.name", "zone", fmt.Sprintf("%s.", rName)),
					resource.TestCheckResourceAttr("gravity_dns_zone.name", "authoritative", "true"),
					resource.TestCheckResourceAttr("gravity_dns_zone.name", "handlers.#", "1"),
					resource.TestCheckResourceAttr("gravity_dns_zone.name", "handlers.0.type", "etcd"),
					resource.TestCheckResourceAttr("gravity_dns_record.record", "fqdn", fmt.Sprintf("foo.%s.", rName)),
				),
			},
		},
	})
}

func testAccResourceDNSRecordSimple(name string) string {
	return fmt.Sprintf(`
resource "gravity_dns_zone" "name" {
  zone          = "%[1]s."
  authoritative = true
  handlers      = [
    {
      type = "etcd",
    }
  ]
}

resource "gravity_dns_record" "record" {
  zone = gravity_dns_zone.name.zone
  hostname = "foo"
  uid = "0"
  data = "1.1.1.1"
  type = "A"
}
`, name)
}
