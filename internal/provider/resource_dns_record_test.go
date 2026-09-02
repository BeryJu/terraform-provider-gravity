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
					resource.TestCheckResourceAttr("gravity_dns_zone.name", "name", fmt.Sprintf("%s.", rName)),
					resource.TestCheckResourceAttr("gravity_dns_zone.name", "authoritative", "true"),
					resource.TestCheckResourceAttr("gravity_dns_zone.name", "handler_configs", "[{\"type\":\"etcd\"}]"),
					resource.TestCheckResourceAttr("gravity_dns_record.record", "fqdn", fmt.Sprintf("foo.%s.", rName)),
					resource.TestCheckResourceAttr("gravity_dns_record.soa", "soa_mbox", "hostmaster.example.com."),
					resource.TestCheckResourceAttr("gravity_dns_record.soa", "soa_serial", "2026090201"),
					resource.TestCheckResourceAttr("gravity_dns_record.soa", "soa_refresh", "7200"),
					resource.TestCheckResourceAttr("gravity_dns_record.soa", "soa_retry", "3600"),
					resource.TestCheckResourceAttr("gravity_dns_record.soa", "soa_expire", "1209600"),
				),
			},
			{
				ResourceName:      "gravity_dns_record.record",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "gravity_dns_record.soa",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceDNSRecordSimple(name string) string {
	return fmt.Sprintf(`
resource "gravity_dns_zone" "name" {
  name            = "%[1]s."
  authoritative   = true
  handler_configs = jsonencode([
    {
      type = "etcd",
    }
  ])
}

resource "gravity_dns_record" "record" {
  zone = gravity_dns_zone.name.name
  hostname = "foo"
  uid = "0"
  data = "1.1.1.1"
  type = "A"
}

resource "gravity_dns_record" "soa" {
  zone        = gravity_dns_zone.name.name
  hostname    = "@"
  uid         = "0"
  type        = "SOA"
  data        = "ns1.example.com."
  ttl         = 3600
  soa_mbox    = "hostmaster.example.com."
  soa_serial  = 2026090201
  soa_refresh = 7200
  soa_retry   = 3600
  soa_expire  = 1209600
}
`, name)
}
