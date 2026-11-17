resource "mailcow_syncjob_source" "example" {
  name = "example.com"
  host = "imap.example.com"
  port = 993
  enc  = "SSL"
}
