# Terraform provider mailcow

terraform provider for [mailcow](https://github.com/mailcow/mailcow-dockerized)

## Requirements

* [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk)
* [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs)

## Disclaimer

This is under development. You will certainly find bugs and limitations. 
In those cases, please report issues or, if you can, submit a pull-request.

### To change

* The API of mailcow for tags always [adds the tags instead of replacing](https://github.com/mailcow/mailcow-dockerized/issues/4681). 
  Either the API has to be changed or extra resources  
  mailcow_domain_tags and mailcow_mailbox_tag have to be implemented.

### Blocked

* There is no API to get user-acl, [but a merge request adding this](https://github.com/mailcow/mailcow-dockerized/pull/4690).
  Without a possibility to read the user-acl a terraform-resource user-acl can not be implemented.