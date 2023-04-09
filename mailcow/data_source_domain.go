package mailcow

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"reflect"
)

func dataSourceDomain() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDomainRead,
		Schema: map[string]*schema.Schema{
			"active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"aliases_left": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"aliases": {
				Type:        schema.TypeInt,
				Description: "limit count of aliases associated with this domain",
				Computed:    true,
			},
			"backupmx": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"bytes_total": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"domain_admins": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"defquota": {
				Type:        schema.TypeInt,
				Description: "predefined mailbox quota in add mailbox form",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of domain",
				Computed:    true,
			},
			"domain": {
				Type:        schema.TypeString,
				Description: "Fully qualified domain name",
				Required:    true,
			},
			"gal": {
				Type:        schema.TypeBool,
				Description: "is domain global address list active or not, it enables shared contacts accross domain in SOGo webmail",
				Computed:    true,
			},
			"mailboxes": {
				Type:        schema.TypeInt,
				Description: "limit count of mailboxes associated with this domain",
				Computed:    true,
			},
			"maxquota": {
				Type:        schema.TypeInt,
				Description: "maximum quota per mailbox",
				Computed:    true,
			},
			"quota": {
				Type:        schema.TypeInt,
				Description: "maximum quota for this domain (for all mailboxes in sum)",
				Computed:    true,
			},
			"mboxes_in_domain": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"mboxes_left": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"msgs_total": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"quota_used_in_domain": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"relay_all_recipients": {
				Type:        schema.TypeBool,
				Description: "if not, them you have to create \"dummy\" mailbox for each address to relay",
				Computed:    true,
			},
			"relay_unknown_only": {
				Type:        schema.TypeBool,
				Description: "Relay non-existing mailboxes only. Existing mailboxes will be delivered locally.",
				Computed:    true,
			},
			"rate_limit": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	id := d.Get("domain").(string)

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

	//exclude := []string{"tags"}
	err = setResourceData(dataSourceDomain(), d, &domain, nil, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	domainTags := domain["tags"]
	if domainTags != nil {
		numTags := len(domainTags.([]interface{}))
		tags := make([]string, numTags)
		for i, tag := range domainTags.([]interface{}) {
			tags[i] = tag.(string)
		}
		err = d.Set("tags", tags)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		tags := make([]string, 0, 0)
		err = d.Set("tags", tags)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(id)

	return diags
}
