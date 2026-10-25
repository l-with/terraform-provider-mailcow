resource "mailcow_mailbox" "demo" {
  domain     = "440044.xyz"
  local_part = "test"
  full_name  = "Test"
  password   = "initial secretpassord"
}

# With Terraform 1.11 or later the password can be write-only: it is sent to
# mailcow but never kept in the plan or the state. Bump the version to send a
# new one.
resource "mailcow_mailbox" "write_only" {
  domain              = "440044.xyz"
  local_part          = "test-wo"
  full_name           = "Test write-only"
  password_wo         = var.initial_password
  password_wo_version = 1
}
