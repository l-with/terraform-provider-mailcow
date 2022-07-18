resource "mailcow_mailbox" "demo" {
  domain     = "440044.xyz"
  local_part = "test"
  password   = "initial secretpassord"
}