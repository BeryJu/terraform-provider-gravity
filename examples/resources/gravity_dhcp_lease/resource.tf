resource "gravity_dhcp_scope" "name" {
  name = "my-scope"
  # Use this scope as a fallback when no scope can be found for a request
  default     = true
  subnet_cidr = "10.10.10.0/24"

  ipam = {
    type        = "internal"
    range_start = "10.10.10.100"
    range_end   = "10.10.10.150"
  }

  # Set DHCP Options
  option {
    tag_name = "router"
    value    = "10.10.10.1"
  }

  # DNS Options
  dns {
    # When `zone` is also configured in gravity, DNS records are created automatically
    zone                 = gravity_dns_zone.example.zone
    add_zone_in_hostname = true
  }
}

resource "gravity_dhcp_lease" "record" {
  scope      = gravity_dhcp_scope.name.name
  hostname   = "foo"
  address    = "10.10.10.25"
  identifier = "aa:bb:cc:dd:ee:ff"
}
