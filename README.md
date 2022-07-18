# Terraform provider mailcow

terraform provider for [mailcow](https://github.com/mailcow/mailcow-dockerized)

## Requirements

* [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk)
* [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs)
* [mailcow Go API Client](https://github.com/l-with/mailcow-go)

# Disclaimer

This is a first implementation. You will certainly find bugs and limitations. In those cases, please report issues or, if you can, submit a pull-request.

## To change

* the API of mailcow for tags always adds the tags instead of replacing (https://github.com/mailcow/mailcow-dockerized/issues/4681). Thus either the API has to be changed or extra resources mailcow_domain_tags and mailcow_mailbox_tag have to be implemented   