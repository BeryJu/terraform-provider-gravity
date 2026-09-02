package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Role configurations are singletons shared by the whole instance, so each
// case restores the value it changed in a final step to keep the instance
// usable for the other tests.

func TestAccResourceRoleTSDB(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "gravity_role_tsdb" "test" {
  enabled = true
  expire  = 900
  scrape  = 60
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gravity_role_tsdb.test", "id", "tsdb"),
					resource.TestCheckResourceAttr("gravity_role_tsdb.test", "enabled", "true"),
					resource.TestCheckResourceAttr("gravity_role_tsdb.test", "expire", "900"),
					resource.TestCheckResourceAttr("gravity_role_tsdb.test", "scrape", "60"),
				),
			},
			{
				ResourceName:      "gravity_role_tsdb.test",
				ImportState:       true,
				ImportStateId:     "tsdb",
				ImportStateVerify: true,
			},
			{
				// Attributes left out fall back to the server's own defaults.
				Config: `
resource "gravity_role_tsdb" "test" {}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gravity_role_tsdb.test", "enabled", "true"),
					resource.TestCheckResourceAttr("gravity_role_tsdb.test", "expire", "1800"),
					resource.TestCheckResourceAttr("gravity_role_tsdb.test", "scrape", "30"),
				),
			},
		},
	})
}

func TestAccResourceRoleBackup(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "gravity_role_backup" "test" {
  endpoint   = "s3.example.com"
  bucket     = "gravity-backups"
  path       = "cluster-a"
  access_key = "an-access-key"
  secret_key = "a-secret-key"
  cron_expr  = "0 */6 * * *"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gravity_role_backup.test", "endpoint", "s3.example.com"),
					resource.TestCheckResourceAttr("gravity_role_backup.test", "bucket", "gravity-backups"),
					resource.TestCheckResourceAttr("gravity_role_backup.test", "path", "cluster-a"),
					resource.TestCheckResourceAttr("gravity_role_backup.test", "access_key", "an-access-key"),
					resource.TestCheckResourceAttr("gravity_role_backup.test", "secret_key", "a-secret-key"),
					resource.TestCheckResourceAttr("gravity_role_backup.test", "cron_expr", "0 */6 * * *"),
				),
			},
			{
				Config: `
resource "gravity_role_backup" "test" {}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gravity_role_backup.test", "endpoint", ""),
					resource.TestCheckResourceAttr("gravity_role_backup.test", "cron_expr", "0 */24 * * *"),
				),
			},
		},
	})
}

func TestAccResourceRoleDNS(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "gravity_role_dns" "test" {
  port = 5353
}
`,
				Check: resource.TestCheckResourceAttr("gravity_role_dns.test", "port", "5353"),
			},
			{
				Config: `
resource "gravity_role_dns" "test" {
  port = 53
}
`,
				Check: resource.TestCheckResourceAttr("gravity_role_dns.test", "port", "53"),
			},
		},
	})
}

// Writing the DHCP role config restarts the DHCP role, and on Gravity 0.33.0
// the role goes on using the context that restart cancelled. Every later write
// to a DHCP scope then fails with `invalid argument: context canceled` for as
// long as the process lives, so running this leaves the instance unable to
// serve the DHCP scope tests. It is opt-in, and needs an instance of its own.
func TestAccResourceRoleDHCP(t *testing.T) {
	if os.Getenv("GRAVITY_TEST_DHCP_ROLE_CONFIG") == "" {
		t.Skip("set GRAVITY_TEST_DHCP_ROLE_CONFIG=1 to run this against a disposable instance")
	}
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "gravity_role_dhcp" "test" {
  port                    = 6767
  lease_negotiate_timeout = 15
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gravity_role_dhcp.test", "port", "6767"),
					resource.TestCheckResourceAttr("gravity_role_dhcp.test", "lease_negotiate_timeout", "15"),
				),
			},
			{
				Config: `
resource "gravity_role_dhcp" "test" {}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gravity_role_dhcp.test", "port", "67"),
					resource.TestCheckResourceAttr("gravity_role_dhcp.test", "lease_negotiate_timeout", "30"),
				),
			},
		},
	})
}

func TestAccResourceRoleTFTP(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "gravity_role_tftp" "test" {
  port         = 6969
  enable_local = true
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gravity_role_tftp.test", "port", "6969"),
					resource.TestCheckResourceAttr("gravity_role_tftp.test", "enable_local", "true"),
				),
			},
			{
				Config: `
resource "gravity_role_tftp" "test" {}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gravity_role_tftp.test", "port", "69"),
					resource.TestCheckResourceAttr("gravity_role_tftp.test", "enable_local", "false"),
				),
			},
		},
	})
}

func TestAccResourceRoleDiscovery(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "gravity_role_discovery" "test" {
  enabled = false
}
`,
				Check: resource.TestCheckResourceAttr("gravity_role_discovery.test", "enabled", "false"),
			},
			{
				Config: `
resource "gravity_role_discovery" "test" {}
`,
				Check: resource.TestCheckResourceAttr("gravity_role_discovery.test", "enabled", "true"),
			},
		},
	})
}

func TestAccResourceRoleAPI(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				// Deliberately leaves `port` at its default: changing it would
				// move the endpoint this provider is talking to.
				//
				// This case also leaves `oidc` unset. Applying an OIDC issuer
				// that is not immediately reachable puts the server into a
				// retry loop, and removing the block while that loop is still
				// running crashes it (a nil dereference in the API role, see
				// pkg/roles/api/auth/method_oidc.go). The provider's side of
				// the OIDC block is covered by TestRoleAPIOIDCRoundTrip.
				Config: `
resource "gravity_role_api" "test" {
  session_duration = "168h"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gravity_role_api.test", "port", "8008"),
					resource.TestCheckResourceAttr("gravity_role_api.test", "session_duration", "168h"),
					resource.TestCheckResourceAttrSet("gravity_role_api.test", "cookie_secret"),
					resource.TestCheckResourceAttr("gravity_role_api.test", "oidc.#", "0"),
				),
			},
			{
				Config: `
resource "gravity_role_api" "test" {}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gravity_role_api.test", "session_duration", ""),
				),
			},
		},
	})
}
