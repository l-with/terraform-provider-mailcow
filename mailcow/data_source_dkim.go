package mailcow

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDkim() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDkimRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pubkey": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"length": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"dkim_txt": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dkim_selector": {
				Type:     schema.TypeString,
				Computed: true,
			},
			//"privkey": {
			//	Type:      schema.TypeString,
			//	Computed:  true,
			//	Sensitive: true,
			//},
		},
	}
}

func dataSourceDkimRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	domain := d.Get("domain").(string)

	request := c.client.Api.MailcowGetDkim(ctx, domain)

	dkim, err := readRequest(request)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(dkim) == 0 {
		return diag.FromErr(errors.New(fmt.Sprint("dkim for domain '", domain, "' not found")))
	}

	dkim["domain"] = domain
	for key, elem := range dataSourceDkim().Schema {
		err = resourceDataSet(d, key, dkim[key], elem)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(domain)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
