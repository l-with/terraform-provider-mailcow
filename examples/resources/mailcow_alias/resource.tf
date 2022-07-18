# create alias
resource "mailcow_alias" "demo" {
  address = "alias-demo@440044.xyz"
  goto    = "demo@440044.xyz"
}
