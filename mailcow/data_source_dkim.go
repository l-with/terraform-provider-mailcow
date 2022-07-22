package mailcow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io"
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
				Type:     schema.TypeString,
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

	request := c.client.DKIMApi.GetDKIMKey(ctx, domain)

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

	dkim := make(map[string]interface{}, 0)
	err = json.NewDecoder(response.Body).Decode(&dkim)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(dkim) == 0 {
		return diag.FromErr(errors.New(fmt.Sprint("dkim for domain '", domain, "' not found")))
	}

	for _, argument := range []string{
		"pubkey",
		"length",
		"dkim_txt",
		"dkim_selector",
		//"privkey",
	} {
		dkimArgument := dkim[argument]
		if dkimArgument != nil {
			err = d.Set(argument, dkimArgument)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(domain)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
