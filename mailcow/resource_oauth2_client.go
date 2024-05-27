package mailcow

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/l-with/terraform-provider-mailcow/api"
	"log"
)

func resourceOAuth2Client() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOAuth2ClientCreate,
		ReadContext:   resourceOAuth2ClientRead,
		DeleteContext: resourceOAuth2ClientDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceOAuth2ClientImport,
		},

		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"redirect_uri": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scope": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOAuth2ClientImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

func resourceOAuth2ClientCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	redirectUri := d.Get("redirect_uri").(string)
	log.Print("[TRACE] resourceOAuth2ClientCreate getId: ", redirectUri)
	id, err := getId(ctx, c.client, redirectUri)
	if err == nil {
		log.Print("[TRACE] resourceOAuth2ClientCreate getId: ", redirectUri, " => ", *id)
		log.Print("[TRACE] resourceOAuth2ClientCreate id: ", *id)
		return diag.Errorf("redirect_uri exists: %s", redirectUri)
	}
	log.Print("[TRACE] resourceOAuth2ClientCreate getId: ", redirectUri, " => error")

	mailcowCreateRequest := api.NewCreateOAuth2ClientRequest()

	err = mailcowCreate(ctx, resourceOAuth2Client(), d, redirectUri, nil, nil, mailcowCreateRequest, c)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err = getId(ctx, c.client, redirectUri)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(*id)

	return resourceOAuth2ClientRead(ctx, d, m)
}

func getId(ctx context.Context, client *api.APIClient, redirectUri string) (*string, error) {
	request := client.Api.MailcowGetOAuth2Clients(ctx)

	oAuth2Clients, err := readAllRequest(request)
	if err != nil {
		return nil, err
	}

	for _, oAuth2Client := range oAuth2Clients {
		if oAuth2Client["redirect_uri"] == redirectUri {
			id := fmt.Sprint(oAuth2Client["id"])
			return &id, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("redirect_uri not found: %s", redirectUri))
}

func resourceOAuth2ClientRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*APIClient)
	id := d.Id()

	request := c.client.Api.MailcowGetOAuth2Client(ctx, id)

	oAuth2Client, err := readRequest(request)
	if err != nil {
		return diag.FromErr(err)
	}

	err = setResourceData(resourceOAuth2Client(), d, &oAuth2Client, nil, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	return diags
}

func resourceOAuth2ClientDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	mailcowDeleteRequest := api.NewDeleteOAuth2ClientRequest()
	diags, _ := mailcowDelete(ctx, d, mailcowDeleteRequest, c)
	return diags
}
