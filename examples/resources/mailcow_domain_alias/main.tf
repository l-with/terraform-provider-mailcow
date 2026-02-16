resource "mailcow_domain_alias" "example" {
  alias_domain  = "alias-domain.tld"
  target_domain = "target-domain.tld"
}
