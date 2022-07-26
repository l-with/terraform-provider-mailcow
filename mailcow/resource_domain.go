package mailcow

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/l-with/terraform-provider-mailcow/api"
	"reflect"
	"strconv"
)

func resourceDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainCreate,
		ReadContext:   resourceDomainRead,
		UpdateContext: resourceDomainUpdate,
		DeleteContext: resourceDomainDelete,
		Schema: map[string]*schema.Schema{
			"active": {
				Type:        schema.TypeBool,
				Description: "is domain active or not",
				Default:     true,
				Optional:    true,
			},
			"aliases": {
				Type:        schema.TypeInt,
				Description: "limit count of aliases associated with this domain",
				Default:     400,
				Optional:    true,
			},
			"backupmx": {
				Type:        schema.TypeBool,
				Description: "relay domain or not",
				Default:     false,
				Optional:    true,
			},
			"defquota": {
				Type:        schema.TypeInt,
				Description: "predefined mailbox quota in add mailbox form",
				Default:     3072,
				Optional:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of domain",
				Default:     "mailcow domain",
				Optional:    true,
			},
			"domain": {
				Type:        schema.TypeString,
				Description: "Fully qualified domain name",
				Required:    true,
				ForceNew:    true,
			},
			"gal": {
				Type:        schema.TypeBool,
				Description: "is domain global address list active or not, it enables shared contacts accross domain in SOGo webmail",
				Default:     false,
				Optional:    true,
			},
			"mailboxes": {
				Type:        schema.TypeInt,
				Description: "limit count of mailboxes associated with this domain",
				Default:     10,
				Optional:    true,
			},
			"maxquota": {
				Type:        schema.TypeInt,
				Description: "maximum quota per mailbox",
				Default:     10240,
				Optional:    true,
			},
			"quota": {
				Type:        schema.TypeInt,
				Description: "maximum quota for this domain (for all mailboxes in sum)",
				Default:     10240,
				Optional:    true,
			},
			"relay_all_recipients": {
				Type:        schema.TypeBool,
				Description: "if not, them you have to create \"dummy\" mailbox for each address to relay",
				Default:     false,
				Optional:    true,
			},
			"relay_unknown_only": {
				Type:        schema.TypeBool,
				Description: "Relay non-existing mailboxes only. Existing mailboxes will be delivered locally.",
				Default:     false,
				Optional:    true,
			},
			"rate_limit": {
				Type:        schema.TypeString,
				Description: "rate limit, decimal with unit s,m,h,d",
				Default:     "10s",
				Optional:    true,
			},
		},
	}
}

func resourceDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*APIClient)

	mailcowCreateRequest := api.NewCreateDomainRequest()

	exclude := []string{"rate_limit"}
	value, ok := d.GetOk("rate_limit")
	if ok {
		rateLimit := value.(string)
		rateFrame := rateLimit[len(rateLimit)-1:]
		mailcowCreateRequest.Set("rl_frame", rateFrame)
		rateValue, err := strconv.Atoi(rateLimit[0 : len(rateLimit)-1])
		if err != nil {
			return diag.FromErr(err)
		}
		mailcowCreateRequest.Set("rl_value", float32(rateValue))

	}
	createRequestSet(mailcowCreateRequest, resourceDomain(), d, &exclude, nil)

	domain := d.Get("domain").(string)
	err := mailcowCreate(ctx, resourceDomain(), d, domain, &exclude, nil, mailcowCreateRequest, c)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(domain)
	return diags
}

func resourceDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*APIClient)
	id := d.Id()

	request := c.client.Api.MailcowGetDomain(ctx, id)

	domain, err := readRequest(request)
	if err != nil {
		return diag.FromErr(err)
	}

	if domain["domain_name"] == nil {
		return diag.FromErr(errors.New("domain not found: " + id))
	}

	domain["aliases"] = domain["max_num_aliases_for_domain"]
	domain["defquota"] = int(domain["def_new_mailbox_quota"].(float64)) / (1024 * 1024)
	domain["domain"] = domain["domain_name"]
	domain["maxquota"] = int(domain["max_quota_for_mbox"].(float64)) / (1024 * 1024)
	domain["mailboxes"] = domain["max_num_mboxes_for_domain"]
	domain["maxquota"] = int(domain["max_quota_for_mbox"].(float64)) / (1024 * 1024)
	domain["quota"] = int(domain["max_quota_for_domain"].(float64)) / (1024 * 1024)
	domainRl := domain["rl"]
	if domainRl != nil {
		if reflect.ValueOf(domainRl).Kind() != reflect.Bool {
			rl := make(map[string]string)
			value := reflect.ValueOf(domainRl)
			for _, key := range value.MapKeys() {
				rl[fmt.Sprint(key)] = fmt.Sprint(value.MapIndex(key))
			}
			domain["rate_limit"] = rl["value"] + rl["frame"]
		}
	}

	setResourceData(resourceDomain(), d, &domain, nil, nil)

	d.SetId(id)

	return diags
}

func resourceDomainUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	mailcowUpdateRequest := api.NewUpdateDomainRequest()

	if d.HasChange("rate_limit") {
		iRateLimit := d.Get("rate_limit")
		rateLimit := iRateLimit.(string)
		rateFrame := rateLimit[len(rateLimit)-1:]
		mailcowUpdateRequest.SetAttr("rl_frame", rateFrame)
		rateValue, err := strconv.Atoi(rateLimit[0 : len(rateLimit)-1])
		if err != nil {
			return diag.FromErr(err)
		}
		mailcowUpdateRequest.SetAttr("rl_value", float32(rateValue))
	}

	updateExclude := []string{
		"rate_limit",
	}
	err := mailcowUpdate(ctx, resourceDomain(), d, &updateExclude, nil, mailcowUpdateRequest, c)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDomainRead(ctx, d, m)
}

func resourceDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	mailcowDeleteRequest := api.NewDeleteDomainRequest()
	diags, _ := mailcowDelete(ctx, d, mailcowDeleteRequest, c)
	return diags
}
