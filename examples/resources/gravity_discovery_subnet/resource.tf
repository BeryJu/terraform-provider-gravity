resource "gravity_discovery_subnet" "office" {
  name        = "office"
  subnet_cidr = "10.10.0.0/24"
  # Resolver used to look up hostnames for the devices that are found
  dns_resolver = "10.10.0.1:53"
  # Forget devices that have not been seen for a day
  discovery_ttl = 86400
}
