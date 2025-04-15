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

## Testing

To make all tests pass, an IdP must be configured. Set a sample config on your test environmnent:

```bash
curl -k 'https://demo.mailcow.email/api/v1/edit/identity-provider' \
  -X POST \
  -H "X-API-Key: 390448-22B69F-FA37D9-19701B-6F033F" \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' --data \
  '{
    "items": ["identity-provider"],
    "attr": {
      "authsource": "keycloak",
      "server_url": "https://auth.demo.mailcow.tld",
      "realm": "mailcow",
      "client_id": "mailcow_terraform",
      "client_secret": "example",
      "redirect_url": "https://demo.mailcow.email",
      "version": "26.1.3"
    }
  }'
```

If no IdP is configured, the mailbox tests will fail with the following message:

```code
resource_mailbox_test.go:18: Step 4/7 error: Check failed: Check 1/1 error: mailcow_mailbox.mailbox: Attribute 'authsource' expected "keycloak", got "mailcow"
```

Then execute tests:

```bash
MAILCOW_HOST_NAME=demo.mailcow.email MAILCOW_API_KEY="390448-22B69F-FA37D9-19701B-6F033F" MAILCOW_INSECURE=true make testacc
```
