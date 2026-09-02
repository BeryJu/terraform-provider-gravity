# Text file, served to every client

resource "gravity_tftp_file" "boot_script" {
  host    = "*"
  name    = "boot.ipxe"
  content = <<-EOT
    #!ipxe
    chain http://boot.my-domain.com/menu.ipxe
  EOT
}

# Binary file, served only to clients connecting to a specific address

resource "gravity_tftp_file" "undionly" {
  host           = "10.10.0.1"
  name           = "undionly.kpxe"
  content_base64 = filebase64("${path.module}/undionly.kpxe")
}
