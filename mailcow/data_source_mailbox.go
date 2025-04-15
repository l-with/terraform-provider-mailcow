package mailcow

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMailbox() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMailboxRead,
		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Description: "e-mail address",
				Required:    true,
			},
			"active": {
				Type:        schema.TypeBool,
				Description: "is alias active or not",
				Computed:    true,
			},
			"domain": {
				Type:        schema.TypeString,
				Description: "domain name",
				Computed:    true,
			},
			"local_part": {
				Type:        schema.TypeString,
				Description: "left part of email address",
				Computed:    true,
			},
			"authsource": {
				Type:        schema.TypeString,
				Description: "Authentication source",
				Computed:    true,
			},
			"full_name": {
				Type:        schema.TypeString,
				Description: "Full name of the mailbox user",
				Computed:    true,
			},
			"quota": {
				Type:        schema.TypeInt,
				Description: "mailbox quota",
				Computed:    true,
			},
			"force_pw_update": {
				Type:        schema.TypeBool,
				Description: "forces the user to update its password on first login",
				Computed:    true,
			},
			"tls_enforce_in": {
				Type:        schema.TypeBool,
				Description: "force inbound email tls encryption",
				Computed:    true,
			},
			"tls_enforce_out": {
				Type:        schema.TypeBool,
				Description: "force outbound tmail tls encryption",
				Computed:    true,
			},
			"sogo_access": {
				Type:        schema.TypeBool,
				Description: "if direct login access to SOGo is granted",
				Computed:    true,
			},
			"imap_access": {
				Type:        schema.TypeBool,
				Description: "if 'IMAP' is an allowed protocol",
				Computed:    true,
			},
			"pop3_access": {
				Type:        schema.TypeBool,
				Description: "if 'POP3' is an allowed protocol",
				Computed:    true,
			},
			"smtp_access": {
				Type:        schema.TypeBool,
				Description: "if 'SMTP' is an allowed protocol",
				Computed:    true,
			},
			"sieve_access": {
				Type:        schema.TypeBool,
				Description: "if 'Sieve' is an allowed protocol",
				Computed:    true,
			},
			//"relayhost": "0",
			//"passwd_update": "2022-07-15 20:31:51",
			//"mailbox_format": "maildir:",
			//"quarantine_notification": "hourly",
			//"quarantine_category": "reject"
		},
	}
}

func dataSourceMailboxRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	id := d.Get("address").(string)

	request := c.client.Api.MailcowGetMailbox(ctx, id)

	mailbox, err := readRequest(request)
	if err != nil {
		return diag.FromErr(err)
	}

	if mailboxAddress(mailbox) != id {
		return diag.FromErr(errors.New(fmt.Sprint("mailbox '", id, "' not found")))
	}

	exclude := []string{"password"}
	mailboxAttributes := []string{
		"force_pw_update",
		"tls_enforce_in",
		"tls_enforce_out",
		"sogo_access",
		"imap_access",
		"pop3_access",
		"smtp_access",
		"sieve_access",
	}
	mailbox["address"] = id
	mailbox["full_name"] = mailbox["name"]
	excludeAndAttributes := append(exclude, mailboxAttributes...)
	err = setResourceData(resourceMailbox(), d, &mailbox, &excludeAndAttributes, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	attributes := mailbox["attributes"].(map[string]interface{})
	err = setResourceData(resourceMailbox(), d, &attributes, &exclude, &mailboxAttributes)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	return diags
}

func mailboxAddress(mailbox map[string]interface{}) string {
	localPart := ""
	mailboxLocalPart := mailbox["local_part"]
	if mailboxLocalPart != nil {
		localPart = mailboxLocalPart.(string)
	}
	domain := ""
	mailboxDomain := mailbox["domain"]
	if mailboxDomain != nil {
		domain = mailboxDomain.(string)
	}
	return localPart + "@" + domain
}
