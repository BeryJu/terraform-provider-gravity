# Authoritative zone

resource "gravity_dns_zone" "example" {
  # Make sure zone ends with a trailing slash
  name          = "my-domain.com."
  authoritative = true
  handlers = [
    {
      type = "memory",
    },
    {
      type = "etcd",
    },
  ]
}

# Forwarding zone

resource "gravity_dns_zone" "forward" {
  # Root zone, will be used for all queries that don't match other zones
  name = "."
  handlers = jsonencode([
    {
      type = "memory",
    },
    {
      type = "etcd",
    },
    {
      type = "forward_ip",
      to   = ["1.1.1.1", "8.8.8.8"],
      # Cache queries and their responses in etcd for 3600s
      cache_ttl = 3600
    },
  ])
}
