---
page_title: "mailcow_dkim Resource - terraform-provider-mailcow"
subcategory: ""
description: |-
---

# mailcow_dkim (Resource)

Provides a DKIM for a domain in mailcow. This can be used to create and delete DKIM for domains.

## Example Usage
```terraform
resource "mailcow_domain" "demo" {
  domain = "440044.xyz"
}

resource "mailcow_dkim" "dkim" {
  domain = mailcow_domain.demo.id
  length = 2048
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `domain` (String)
- `length` (Number)

### Optional

- `dkim_selector` (String)

### Read-Only

- `dkim_txt` (String)
- `id` (String) The ID of this resource.
- `pubkey` (String)
