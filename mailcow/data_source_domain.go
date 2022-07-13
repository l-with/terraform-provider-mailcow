package mailcow

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io"
	"reflect"
	"strconv"
)

func dataSourceDomain() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDomainRead,
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			//"active_int": {
			//	Type:     schema.TypeInt,
			//	Computed: true,
			//},
			"aliases_left": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"aliases_in_domain": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"backupmx": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			//"backupmx_int": {
			//	Type:     schema.TypeInt,
			//	Computed: true,
			//},
			"bytes_total": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"domain_admins": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"def_quota_for_mbox": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"def_new_mailbox_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gal": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			//"gal_int": {
			//	Type:     schema.TypeInt,
			//	Computed: true,
			//},
			"max_new_mailbox_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_num_aliases_for_domain": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_num_mboxes_for_domain": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_quota_for_domain": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_quota_for_mbox": {
				Type:     schema.TypeInt,
				Computed: true,
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
				Type:     schema.TypeBool,
				Computed: true,
			},
			//"relay_all_recipients_int": {
			//	Type:     schema.TypeInt,
			//	Computed: true,
			//},
			"relayhost": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"relay_unknown_only": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			//"relay_unknown_only_int": {
			//	Type:     schema.TypeInt,
			//	Computed: true,
			//},
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
	domainName := d.Get("domain_name").(string)

	request := c.client.DomainsApi.GetDomains(ctx, domainName)

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

	domain := make(map[string]interface{}, 0)
	err = json.NewDecoder(response.Body).Decode(&domain)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, argument := range []string{
		//"active_int",
		"aliases_left",
		"aliases_left",
		"aliases_in_domain",
		//"backupmx_int",
		"bytes_total",
		"domain_admins",
		"def_quota_for_mbox",
		"def_new_mailbox_quota",
		"description",
		//"gal_int",
		"max_new_mailbox_quota",
		"max_num_aliases_for_domain",
		"max_num_mboxes_for_domain",
		"max_quota_for_domain",
		"max_quota_for_mbox",
		"mboxes_in_domain",
		"mboxes_left",
		"msgs_total",
		//"relay_all_recipients_int",
		//"relay_unknown_only_int",
		"relayhost",
	} {
		err = d.Set(argument, domain[argument])
		if err != nil {
			return diag.FromErr(err)
		}
	}

	for _, argument := range []string{
		"active",
		"backupmx",
		"gal",
		"relay_all_recipients",
		"relay_unknown_only",
	} {
		boolValue := false
		if int(domain[argument].(float64)) >= 1 {
			boolValue = true
		}
		err = d.Set(argument, boolValue)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	quotaUsedInDomain, err := strconv.Atoi(fmt.Sprint(domain["quota_used_in_domain"]))
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("quota_used_in_domain", quotaUsedInDomain)
	if err != nil {
		return diag.FromErr(err)
	}

	rl := make(map[string]string)
	value := reflect.ValueOf(domain["rl"])
	for _, key := range value.MapKeys() {
		rl[fmt.Sprint(key)] = fmt.Sprint(value.MapIndex(key))
	}
	rateLimit := rl["value"] + rl["frame"]

	err = d.Set("rate_limit", rateLimit)
	if err != nil {
		return diag.FromErr(err)
	}

	if domain["tags"] != nil {
		numTags := len(domain["tags"].([]interface{}))
		tags := make([]string, numTags)
		for i, tag := range domain["tags"].([]interface{}) {
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

	d.SetId(domainName)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
