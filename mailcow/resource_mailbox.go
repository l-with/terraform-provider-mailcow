package mailcow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/l-with/terraform-provider-mailcow/api"
)

const (
	mailcowAuthsourceInternal = "mailcow"
	mailcowAuthsourceKeycloak = "keycloak"
	mailcowAuthsourceLdap     = "ldap"
	mailcowAuthsourceOidc     = "generic-oidc"
)

func resourceMailbox() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMailboxCreate,
		ReadContext:   resourceMailboxRead,
		UpdateContext: resourceMailboxUpdate,
		DeleteContext: resourceMailboxDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceMailboxImport,
		},

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
			"authsource": {
				Type:         schema.TypeString,
				Description:  "Authentication source to use. One of: generic-oidc, mailcow, keycloak, ldap.",
				Default:      mailcowAuthsourceInternal,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{mailcowAuthsourceInternal, mailcowAuthsourceKeycloak, mailcowAuthsourceLdap, mailcowAuthsourceOidc}, false),
			},
			"address": {
				Type:        schema.TypeString,
				Description: "e-mail address",
				Computed:    true,
			},
			"full_name": {
				Type:        schema.TypeString,
				Description: "Full name of the mailbox user",
				Required:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "mailbox password (the password is excluded from update)",
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
				Description: "force outbound mail tls encryption",
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

func resourceMailboxImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

func resourceMailboxCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*APIClient)

	mailcowCreateRequest := api.NewCreateMailboxRequest()

	address := d.Get("local_part").(string) + "@" + d.Get("domain").(string)
	err := d.Set("address", address)
	if err != nil {
		return diag.FromErr(err)
	}

	mailcowCreateRequest.Set("password2", d.Get("password"))

	mapArguments := map[string]string{"full_name": "name"}

	err = mailcowCreate(ctx, resourceMailbox(), d, address, nil, &mapArguments, mailcowCreateRequest, c)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(address)

	return diags
}

func resourceMailboxRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*APIClient)
	id := d.Id()

	request := c.client.Api.MailcowGetMailbox(ctx, id)

	mailbox, err := readRequest(request)
	if err != nil {
		return diag.FromErr(err)
	}

	exclude := []string{
		"password",
	}
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
	mailbox["quota"] = int(mailbox["quota"].(float64)) / (1024 * 1024)

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

func resourceMailboxUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	mailcowUpdateRequest := api.NewUpdateMailboxRequest()

	exclude := []string{
		"password",
	}
	mapArguments := map[string]string{
		"full_name": "name",
	}
	err := mailcowUpdate(ctx, resourceMailbox(), d, &exclude, &mapArguments, mailcowUpdateRequest, c)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceMailboxRead(ctx, d, m)
}

func resourceMailboxDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	mailcowDeleteRequest := api.NewDeleteMailboxRequest()
	diags, _ := mailcowDelete(ctx, d, mailcowDeleteRequest, c)
	return diags
}
