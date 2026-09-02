resource "gravity_role_tsdb" "tsdb" {
  enabled = true
  # Keep metrics for 30 minutes, collecting every 30 seconds
  expire = 1800
  scrape = 30
}
