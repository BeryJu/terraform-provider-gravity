resource "gravity_dns_zone" "example" {
  # Make sure zone ends with a trailing slash
  zone          = "my-domain.com."
  authoritative = true
  handler_configs = jsonencode([
    {
      type = "memory",
    },
    {
      type = "etcd",
    },
  ])
}

resource "gravity_dns_record" "record" {
  zone     = gravity_dns_zone.example.zone
  hostname = "foo"
  uid      = "0"
  data     = "1.1.1.1"
  type     = "A"
}
