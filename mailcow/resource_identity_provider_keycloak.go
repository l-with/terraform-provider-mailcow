package mailcow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/l-with/terraform-provider-mailcow/api"
)

func resourceIdentityProviderKeycloak() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityProviderKeycloakCreate,
		ReadContext:   resourceIdentityProviderKeycloakRead,
		DeleteContext: resourceIdentityProviderKeycloakDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceIdentityProviderKeycloakImport,
		},

		Schema: map[string]*schema.Schema{
			"authsource": {
				Type:        schema.TypeString,
				Description: "must be 'keycloak'",
				Optional:    true,
				Default:     "keycloak",
				ForceNew:    true,
			},
			"client_id": {
				Type:        schema.TypeString,
				Description: "the Client ID assigned to mailcow Client in Keycloak",
				Required:    true,
				ForceNew:    true,
			},
			"client_secret": {
				Type:        schema.TypeString,
				Description: "the Client Secret assigned to the mailcow client in Keycloak",
				Required:    true,
				ForceNew:    true,
			},
			"import_users": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"ignore_ssl_error": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"mailpassword_flow": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"periodic_sync": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},
			"realm": {
				Type:        schema.TypeString,
				Description: "the Keycloak realm where the mailcow client is configured",
				Required:    true,
				ForceNew:    true,
			},
			"redirect_url": {
				Type:        schema.TypeString,
				Description: "the redirect URL that Keycloak will use after authentication. This should point to your mailcow UI. Example: https://mail.mailcow.tld",
				Required:    true,
				ForceNew:    true,
			},
			"server_url": {
				Type:        schema.TypeString,
				Description: "the base URL of the Keycloak server",
				Required:    true,
				ForceNew:    true,
			},
			"sync_interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  15,
				ForceNew: true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "specifies the Keycloak version (cite from blog 'It is essential to know whether a version greater or smaller than 20 is used since mailcow needs to add the '“openid'” scope accordingly.')",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceIdentityProviderKeycloakImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

func resourceIdentityProviderKeycloakCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	mailcowUpdateRequest := api.NewUpdateIdentityProviderKeycloakRequest()

	err := mailcowUpdate(ctx, resourceIdentityProviderKeycloak(), d, nil, nil, mailcowUpdateRequest, c)
	if err != nil {
		return diag.FromErr(err)
	}
	ignore_ssl_error := d.Get("ignore_ssl_error")
	if ignore_ssl_error == "" {
		d.Set("ignore_ssl_error", false)
	}

	return resourceIdentityProviderKeycloakRead(ctx, d, m)
}

func resourceIdentityProviderKeycloakRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*APIClient)

	request := c.client.Api.MailcowGetIdentityProviderKeycloak(ctx)

	identityProviderKeycloak, err := readRequest(request)
	if err != nil {
		return diag.FromErr(err)
	}

	err = setResourceData(resourceIdentityProviderKeycloak(), d, &identityProviderKeycloak, nil, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	authsource := d.Get("authsource")
	d.SetId(authsource.(string))

	return diags
}

func resourceIdentityProviderKeycloakDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	mailcowDeleteRequest := api.NewDeleteDkimRequest()
	diags, _ := mailcowDelete(ctx, d, mailcowDeleteRequest, c)
	return diags
}
