resource "mailcow_mailbox" "demo" {
  domain     = "440044.xyz"
  local_part = "test"
  full_name  = "Test"
  password   = "initial secretpassord"
}

# With Terraform 1.11+ or OpenTofu 1.11+ the password can be write-only: it is
# sent to mailcow but never kept in the plan or the state. An ephemeral variable
# keeps it out of a saved plan as well; pass it with TF_VAR_mailbox_password
# rather than a .tfvars file.
variable "mailbox_password" {
  type      = string
  sensitive = true
  ephemeral = true
}

resource "mailcow_mailbox" "write_only" {
  domain      = "440044.xyz"
  local_part  = "test-wo"
  full_name   = "Test write-only"
  password_wo = var.mailbox_password

  # Optional. Without it the password is set on create and never touched
  # again, as with `password`. Bump it to send a new password.
  password_wo_version = 1
}
