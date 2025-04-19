resource "mailcow_identity_provider_keycloak" "keycloak" {
  authsource    = "keycloak"
  server_url    = "https://auth.demo.mailcow.tld"
  realm         = "mailcow"
  client_id     = "mailcow_terraform"
  client_secret = "example"
  redirect_url  = "https://demo.mailcow.tld"
  version       = "26.1.3"
  import_users  = true
  periodic_sync = true
  sync_interval = 20
}