# Simple user (read-only)

resource "gravity_user" "read-only" {
  username = "read-only"
  permissions = jsonencode([
    {
      path    = "/*",
      methods = ["get", "head"]
    }
  ])
}

# Admin user

resource "gravity_user" "admin" {
  username = "admin"
  permissions = jsonencode([
    {
      path    = "/*",
      methods = ["*"]
    }
  ])
}
