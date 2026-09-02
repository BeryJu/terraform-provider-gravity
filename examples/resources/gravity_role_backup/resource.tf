resource "gravity_role_backup" "backup" {
  endpoint   = "s3.my-domain.com"
  bucket     = "gravity-backups"
  path       = "production"
  access_key = var.s3_access_key
  secret_key = var.s3_secret_key
  # Every 6 hours
  cron_expr = "0 */6 * * *"
}
