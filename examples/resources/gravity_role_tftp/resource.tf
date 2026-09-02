resource "gravity_role_tftp" "tftp" {
  port = 69
  # Also serve files from the instance's local filesystem
  enable_local = false
}
