package provider

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceTFTPFile(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	const initial = "#!ipxe\nchain http://example.com/boot\n"
	const updated = "#!ipxe\nexit\n"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTFTPFile(rName, initial),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gravity_tftp_file.text", "name", rName+".ipxe"),
					resource.TestCheckResourceAttr("gravity_tftp_file.text", "host", "*"),
					resource.TestCheckResourceAttr("gravity_tftp_file.text", "content", initial),
					resource.TestCheckResourceAttr("gravity_tftp_file.text", "size_bytes", strconv.Itoa(len(initial))),
					resource.TestCheckResourceAttr("gravity_tftp_file.binary", "content_base64", testTFTPBinary),
					resource.TestCheckResourceAttr("gravity_tftp_file.binary", "size_bytes", "4"),
				),
			},
			{
				// Updating the contents in place must be reflected on read.
				Config: testAccResourceTFTPFile(rName, updated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gravity_tftp_file.text", "content", updated),
					resource.TestCheckResourceAttr("gravity_tftp_file.text", "size_bytes", strconv.Itoa(len(updated))),
				),
			},
			{
				ResourceName:      "gravity_tftp_file.text",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("*/%s.ipxe", rName),
				ImportStateVerify: true,
				// The API cannot say which encoding the config used, so an
				// imported file always comes back as content.
				ImportStateVerifyIgnore: []string{"content"},
			},
		},
	})
}

// Bytes that are not valid UTF-8, to exercise the base64 path.
var testTFTPBinary = base64.StdEncoding.EncodeToString([]byte{0x00, 0x01, 0x02, 0xff})

func testAccResourceTFTPFile(name string, content string) string {
	return fmt.Sprintf(`
resource "gravity_tftp_file" "text" {
  host    = "*"
  name    = "%[1]s.ipxe"
  content = %[2]q
}

resource "gravity_tftp_file" "binary" {
  host           = "*"
  name           = "%[1]s.bin"
  content_base64 = %[3]q
}
`, name, content, testTFTPBinary)
}
