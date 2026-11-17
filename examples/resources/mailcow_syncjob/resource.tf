resource "mailcow_syncjob_source" "example" {
  name = "example.com"
  host = "imap.example.com"
}

resource "mailcow_syncjob" "syncjob" {
  username  = "demo@440044.xyz"
  source_id = mailcow_syncjob_source.example.id
  user1     = "demo@example.com"
  password1 = "secret-password"

  custom_params {
    option = "errorsmax"
    value  = "100"
  }
}
