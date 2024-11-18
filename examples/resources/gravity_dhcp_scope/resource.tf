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
    # Tag name, can be any of
    # - subnet_mask
    # - router
    # - time_server
    # - name_server
    # - domain_name
    # - bootfile
    # - tftp_server
    tag_name = "router"
    value    = "10.10.10.1"
  }
  option {
    tag_name = "name_server"
    # value = "10.10.10.1"
    # Set the value as base64 to allow for binary data
    value64 = [base64encode("10.10.10.1")]
  }
  option {
    # Option tag as integer
    tag   = 43
    value = "10.10.10.2"
  }

  # DNS Options
  dns {
    # When `zone` is also configured in gravity, DNS records are created automatically
    zone                 = gravity_dns_zone.example.id
    add_zone_in_hostname = true
  }
}

resource "gravity_dns_zone" "example" {
  # Make sure zone ends with a trailing slash
  zone          = "my-domain.com."
  authoritative = true
  handlers = jsonencode([
    {
      type = "memory",
    },
    {
      type = "etcd",
    },
  ])
}
