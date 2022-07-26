package mailcow

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/l-with/terraform-provider-mailcow/api"
)

func resourceDkim() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDkimCreate,
		ReadContext:   resourceDkimRead,
		DeleteContext: resourceDkimDelete,
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"pubkey": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"length": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"dkim_txt": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dkim_selector": {
				Type:     schema.TypeString,
				Default:  "dkim",
				Optional: true,
				ForceNew: true,
			},
			//"privkey": {
			//	Type:      schema.TypeString,
			//	Computed:  true,
			//	Sensitive: true,
			//},
		},
	}
}

func resourceDkimCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	mailcowCreateRequest := api.NewCreateDkimRequest()

	mapArguments := map[string]string{
		"length": "key_size",
		"domain": "domains",
	}
	domain := d.Get("domain").(string)
	err := mailcowCreate(ctx, resourceDkim(), d, domain, nil, &mapArguments, mailcowCreateRequest, c)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(domain)

	return resourceDkimRead(ctx, d, m)
}

func resourceDkimRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*APIClient)
	id := d.Id()

	request := c.client.Api.MailcowGetDkim(ctx, id)

	dkim, err := readRequest(request)
	if err != nil {
		return diag.FromErr(err)
	}

	dkim["domain"] = id

	err = setResourceData(resourceDkim(), d, &dkim, nil, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	return diags
}

func resourceDkimDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	mailcowDeleteRequest := api.NewDeleteDkimRequest()
	diags, _ := mailcowDelete(ctx, d, mailcowDeleteRequest, c)
	return diags
}
