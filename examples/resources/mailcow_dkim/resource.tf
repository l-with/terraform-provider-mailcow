resource "mailcow_domain" "demo" {
  domain = "440044.xyz"
}

resource "mailcow_dkim" "dkim" {
  domain = mailcow_domain.demo.id
  length = 2048
}
