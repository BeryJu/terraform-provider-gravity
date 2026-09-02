resource "gravity_user" "automation" {
  username = "automation"
  permissions = jsonencode([
    {
      path    = "/api/v1/dns/*"
      methods = ["GET", "POST", "DELETE"]
    }
  ])
}

resource "gravity_api_token" "automation" {
  username = gravity_user.automation.username
}

output "automation_token" {
  value     = gravity_api_token.automation.key
  sensitive = true
}
