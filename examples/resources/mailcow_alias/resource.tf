# create alias
resource "mailcow_alias" "alias" {
  address = "alias-demo@440044.xyz"
  goto    = "demo@440044.xyz"
}
