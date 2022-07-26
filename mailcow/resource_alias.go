package mailcow

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/l-with/terraform-provider-mailcow/api"
)

func resourceAlias() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAliasCreate,
		ReadContext:   resourceAliasRead,
		UpdateContext: resourceAliasUpdate,
		DeleteContext: resourceAliasDelete,
		Schema: map[string]*schema.Schema{
			"active": {
				Type:        schema.TypeBool,
				Description: "is alias active or not",
				Default:     true,
				Optional:    true,
			},
			"address": {
				Type:        schema.TypeString,
				Description: "alias address, for catchall use \"@domain.tld\"",
				Required:    true,
			},
			"goto": {
				Type:        schema.TypeString,
				Description: "destination address, comma separated\nSpecial values are spam@locahost, ",
				Required:    true,
			},
			"sogo_visible": {
				Type:        schema.TypeBool,
				Description: "visibility as selectable sender in SOGo",
				Default:     false,
				Optional:    true,
			},
			"private_comment": {
				Type:        schema.TypeString,
				Description: "private comment",
				Optional:    true,
			},
			"public_comment": {
				Type:        schema.TypeString,
				Description: "public comment",
				Optional:    true,
			},
		},
	}
}

func resourceAliasCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*APIClient)

	mailcowCreateRequest := api.NewCreateAliasRequest()

	createRequestSet(mailcowCreateRequest, resourceAlias(), d, nil, nil)

	request := c.client.Api.MailcowCreate(ctx).MailcowCreateRequest(*mailcowCreateRequest)
	response, _, err := c.client.Api.MailcowCreateExecute(request)
	if err != nil {
		return diag.FromErr(err)
	}
	err = checkResponse(response, "resourceAliasCreate", d.Get("address").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	err, id := response.GetId()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*id)

	return diags
}

func resourceAliasRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*APIClient)
	id := d.Id()

	request := c.client.Api.MailcowGetAlias(ctx, id)

	alias, err := readRequest(request)
	if err != nil {
		return diag.FromErr(err)
	}

	if alias["id"] == nil {
		return diag.FromErr(errors.New("alias id not found: " + id))
	}

	setResourceData(resourceAlias(), d, &alias, nil, nil)

	d.SetId(fmt.Sprint(alias["id"].(float64)))

	return diags
}

func resourceAliasUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	mailcowUpdateRequest := api.NewUpdateAliasRequest()

	err := mailcowUpdate(ctx, resourceAlias(), d, nil, nil, mailcowUpdateRequest, c)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceAliasRead(ctx, d, m)
}

func resourceAliasDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	mailcowDeleteRequest := api.NewDeleteAliasRequest()
	diags, _ := mailcowDelete(ctx, d, mailcowDeleteRequest, c)
	return diags
}
