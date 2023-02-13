provider "gravity" {
  url = "http://my-gravity.domain.td:8008"
  # Ignore Certificate errors when using HTTPS
  insecure = false
  token    = "my-token"
}
