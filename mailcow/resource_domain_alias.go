package mailcow

import (
	"context"
	"errors"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/l-with/terraform-provider-mailcow/api"
)

func resourceDomainAlias() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainAliasCreate,
		ReadContext:   resourceDomainAliasRead,
		UpdateContext: resourceDomainAliasUpdate,
		DeleteContext: resourceDomainAliasDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDomainAliasImport,
		},

		Schema: map[string]*schema.Schema{
			"active": {
				Type:        schema.TypeBool,
				Description: "is domain alias active or not",
				Default:     true,
				Optional:    true,
			},
			"alias_domain": {
				Type:        schema.TypeString,
				Description: "Alias domain name",
				Required:    true,
				ForceNew:    true,
			},
			"target_domain": {
				Type:        schema.TypeString,
				Description: "Target domain name",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceDomainAliasImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

func resourceDomainAliasCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	log.Print("resourceDomainAliasCreate")

	mailcowCreateRequest := api.NewCreateAliasDomainRequest()
	aliasDomain := d.Get("alias_domain").(string)
	createRequestSet(mailcowCreateRequest, resourceDomainAlias(), d, nil, nil)

	request := c.client.Api.MailcowCreate(ctx).MailcowCreateRequest(*mailcowCreateRequest)
	response, _, err := c.client.Api.MailcowCreateExecute(request)
	if err != nil {
		return diag.FromErr(err)
	}
	err = checkResponse(response, mailcowCreateRequest.ResourceName, aliasDomain)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := response.GetAliasDomainId()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*id)
	return diags
}

func resourceDomainAliasRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	log.Print("resourceDomainAliasRead")

	c := m.(*APIClient)
	id := d.Id()

	request := c.client.Api.MailcowGetAliasDomain(ctx, id)

	aliasDomain, err := readRequest(request)
	if err != nil {
		return diag.FromErr(err)
	}

	if aliasDomain["alias_domain"] == nil {
		return diag.FromErr(errors.New("domain alias not found: " + id))
	}

	err = setResourceData(resourceDomainAlias(), d, &aliasDomain, nil, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	return diags
}

func resourceDomainAliasUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	mailcowUpdateRequest := api.NewUpdateAliasDomainRequest()

	err := mailcowUpdate(ctx, resourceDomainAlias(), d, nil, nil, mailcowUpdateRequest, c)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDomainAliasRead(ctx, d, m)
}

func resourceDomainAliasDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	mailcowDeleteRequest := api.NewDeleteAliasDomainRequest()
	diags, _ := mailcowDelete(ctx, d, mailcowDeleteRequest, c)
	return diags
}
