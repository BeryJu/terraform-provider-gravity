resource "gravity_role_dhcp" "dhcp" {
  port = 67
  # Wait 30s for other instances to respond before offering a lease
  lease_negotiate_timeout = 30
}
