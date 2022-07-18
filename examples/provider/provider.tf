resource "mailcow_domain" "demo" {
  domain = "440044.xyz"
}

resource "mailcow_mailbox" "demo" {
  domain     = mailcow_domain.demo.domain
  local_part = "test"
  password   = "initial secretpassord"
}

resource "mailcow_alias" "demo" {
  address = "alias-demo@440044.xyz"
  goto    = mailcow_mailbox.demo.id
}

resource "mailcow_dkim" "demo" {
  domain = mailcow_domain.demo.domain
  length = 2048
}
