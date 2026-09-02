resource "gravity_role_api" "api" {
  port             = 8008
  session_duration = "168h"

  oidc {
    client_id            = "gravity"
    client_secret        = var.oidc_client_secret
    issuer               = "https://id.my-domain.com/"
    redirect_url         = "https://gravity.my-domain.com/auth/callback"
    scopes               = ["openid", "email", "profile"]
    token_username_field = "preferred_username"
  }
}
