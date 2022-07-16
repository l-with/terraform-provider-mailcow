package mailcow

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "github.com/l-with/mailcow-go"
	"io"
	"log"
	"reflect"
)

func resourceMailbox() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMailboxCreate,
		ReadContext:   resourceMailboxRead,
		UpdateContext: resourceMailboxUpdate,
		DeleteContext: resourceMailboxDelete,
		Schema: map[string]*schema.Schema{
			"active": {
				Type:        schema.TypeBool,
				Description: "is alias active or not",
				Default:     true,
				Optional:    true,
			},
			"domain": {
				Type:        schema.TypeString,
				Description: "domain name",
				Required:    true,
				ForceNew:    true,
			},
			"local_part": {
				Type:        schema.TypeString,
				Description: "left part of email address",
				Required:    true,
				ForceNew:    true,
			},
			"full_name": {
				Type:        schema.TypeString,
				Description: "Full name of the mailbox user",
				Optional:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "mailbox password",
				Required:    true,
				Sensitive:   true,
			},
			"quota": {
				Type:        schema.TypeInt,
				Description: "mailbox quota",
				Optional:    true,
			},
			"force_pw_update": {
				Type:        schema.TypeBool,
				Description: "forces the user to update its password on first login",
				Default:     true,
				Optional:    true,
			},
			"tls_enforce_in": {
				Type:        schema.TypeBool,
				Description: "force inbound email tls encryption",
				Default:     false,
				Optional:    true,
			},
			"tls_enforce_out": {
				Type:        schema.TypeBool,
				Description: "force outbound tmail tls encryption",
				Default:     false,
				Optional:    true,
			},
			"sogo_access": {
				Type:        schema.TypeBool,
				Description: "if direct login access to SOGo is granted",
				Default:     true,
				Optional:    true,
			},
			"imap_access": {
				Type:        schema.TypeBool,
				Description: "if 'IMAP' is an allowed protocol",
				Default:     true,
				Optional:    true,
			},
			"pop3_access": {
				Type:        schema.TypeBool,
				Description: "if 'POP3' is an allowed protocol",
				Default:     true,
				Optional:    true,
			},
			"smtp_access": {
				Type:        schema.TypeBool,
				Description: "if 'SMTP' is an allowed protocol",
				Default:     true,
				Optional:    true,
			},
			"sieve_access": {
				Type:        schema.TypeBool,
				Description: "if 'Sieve' is an allowed protocol",
				Default:     true,
				Optional:    true,
			},
			//"relayhost": "0",
			//"passwd_update": "2022-07-15 20:31:51",
			//"mailbox_format": "maildir:",
			//"quarantine_notification": "hourly",
			//"quarantine_category": "reject"
		},
	}
}

func resourceMailboxCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*APIClient)

	createMailboxRequest := api.NewCreateMailboxRequest()
	createMailboxRequest.SetActive(d.Get("active").(bool))
	domain := d.Get("domain").(string)
	createMailboxRequest.SetDomain(domain)
	localPart := d.Get("local_part").(string)
	createMailboxRequest.SetLocalPart(localPart)
	name, ok := d.GetOk("full_name")
	if ok {
		createMailboxRequest.SetName(name.(string))
	}
	createMailboxRequest.SetPassword(d.Get("password").(string))
	createMailboxRequest.SetPassword2(d.Get("password").(string))
	quota, ok := d.GetOk("quota")
	if ok {
		createMailboxRequest.SetQuota(float32(quota.(int)))
	}
	createMailboxRequest.SetForcePwUpdate(d.Get("force_pw_update").(bool))
	createMailboxRequest.SetTlsEnforceOut(d.Get("tls_enforce_out").(bool))
	createMailboxRequest.SetTlsEnforceIn(d.Get("tls_enforce_in").(bool))
	createMailboxRequest.SetSogoAccess(d.Get("sogo_access").(bool))
	createMailboxRequest.SetImapAccess(d.Get("imap_access").(bool))
	createMailboxRequest.SetPop3Access(d.Get("pop3_access").(bool))
	createMailboxRequest.SetSmtpAccess(d.Get("smtp_access").(bool))
	createMailboxRequest.SetSieveAccess(d.Get("sieve_access").(bool))

	request := c.client.MailboxesApi.CreateMailbox(ctx).CreateMailboxRequest(*createMailboxRequest)
	_, _, err := c.client.MailboxesApi.CreateMailboxExecute(request)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(localPart + "@" + domain)
	return diags
}

func resourceMailboxRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*APIClient)
	emailAddress := d.Id()

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

func resourceMailboxUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	emailAddress := d.Id()

	updateMailboxRequest := api.NewUpdateMailboxRequest()
	updateMailboxRequestAttr := api.NewUpdateMailboxRequestAttr()

	if d.HasChange("active") {
		updateMailboxRequestAttr.SetActive(d.Get("active").(bool))
	}
	log.Print("[TRACE] resourceMailboxUpdate d.Get(full_name): ", d.Get("full_name"))
	if d.HasChange("full_name") {
		updateMailboxRequestAttr.SetName(d.Get("full_name").(string))
	}
	if d.HasChange("quota") {
		updateMailboxRequestAttr.SetQuota(float32(d.Get("quota").(int)))
	}
	if d.HasChange("sogo_access") {
		updateMailboxRequestAttr.SetSogoAccess(d.Get("sogo_access").(bool))
	}

	items := make([]string, 1)
	items[0] = emailAddress

	updateMailboxRequest.SetItems(items)
	updateMailboxRequest.SetAttr(*updateMailboxRequestAttr)
	request := c.client.MailboxesApi.UpdateMailbox(ctx).UpdateMailboxRequest(*updateMailboxRequest)
	_, _, err := c.client.MailboxesApi.UpdateMailboxExecute(request)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceMailboxRead(ctx, d, m)
}

func resourceMailboxDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*APIClient)

	deleteMailboxRequest := api.NewDeleteMailboxRequest()

	items := make([]string, 1)
	items[0] = d.Id()
	deleteMailboxRequest.SetItems(items)

	request := c.client.MailboxesApi.DeleteMailbox(ctx).DeleteMailboxRequest(*deleteMailboxRequest)
	_, _, err := c.client.MailboxesApi.DeleteMailboxExecute(request)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
