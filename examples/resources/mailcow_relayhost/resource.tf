resource "mailcow_relayhost" "relayhost" {
  hostname = "my-smtp-relay.com:2525"
  username = "my-username"
  password = "my-password"
}