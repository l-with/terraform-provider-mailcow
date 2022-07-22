package mailcow

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io"
	"log"
	"reflect"
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
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	c := m.(*APIClient)
	emailAddress := d.Get("address").(string)

	request := c.client.MailboxesApi.GetMailboxes(ctx, emailAddress)

	response, err := request.Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			diag.FromErr(err)
		}
	}(response.Body)

	mailbox := make(map[string]interface{}, 0)
	err = json.NewDecoder(response.Body).Decode(&mailbox)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, argument := range []string{
		"domain",
		"address",
		"local_part",
		"full_name",
	} {
		err = d.Set(argument, mailbox[argument])
		if err != nil {
			return diag.FromErr(err)
		}
		log.Print("[TRACE] resourceMailboxRead mailbox[", argument, "]: ", mailbox[argument])
	}
	mailboxQuota := mailbox["quota"]
	if mailboxQuota != nil {
		err = d.Set("quota", int(mailboxQuota.(float64))/(1024*1024))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	for _, argumentBool := range []string{
		"active",
	} {
		boolValue := false
		if mailbox[argumentBool] != nil {
			if int(mailbox[argumentBool].(float64)) >= 1 {
				boolValue = true
			}
		}
		err = d.Set(argumentBool, boolValue)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	attributes := mailbox["attributes"]
	if attributes != nil {
		value := reflect.ValueOf(attributes)
		for _, argument := range []string{
			"force_pw_update",
			"tls_enforce_in",
			"tls_enforce_out",
			"sogo_access",
			"imap_access",
			"pop3_access",
			"smtp_access",
			"sieve_access",
		} {
			boolValue := false
			if fmt.Sprint(value.MapIndex(reflect.ValueOf(argument))) == "1" {
				boolValue = true
			}
			err = d.Set(argument, boolValue)
			if err != nil {
				return diag.FromErr(err)
			}
		}
		//for _, key := range value.MapKeys() {
		//	boolValue := false
		//	stringValue := fmt.Sprint(value.MapIndex(key))
		//	if stringValue == "1" {
		//		boolValue = true
		//	}
		//	err = d.Set(fmt.Sprint(key), boolValue)
		//}
	}

	d.SetId(emailAddress)

	return diags
}
