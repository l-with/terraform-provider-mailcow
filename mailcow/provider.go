package mailcow

import (
	"context"
	"github.com/l-with/terraform-provider-mailcow/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MAILCOW_HOST_NAME", nil),
				Description: "The name of the mailcow host, can optionally be passed as `MAILCOW_HOST_NAME` environmental variable",
			},
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("MAILCOW_API_KEY", nil),
				Description: "The mailcow API key, can optionally be passed as `MAILCOW_API_KEY` environmental variable",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"mailcow_alias":         resourceAlias(),
			"mailcow_domain":        resourceDomain(),
			"mailcow_mailbox":       resourceMailbox(),
			"mailcow_dkim":          resourceDkim(),
			"mailcow_syncjob":       resourceSyncjob(),
			"mailcow_oauth2_client": resourceOAuth2Client(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"mailcow_domain":  dataSourceDomain(),
			"mailcow_mailbox": dataSourceMailbox(),
			"mailcow_dkim":    dataSourceDkim(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// APIClient Hold the API Client and any relevant configuration
type APIClient struct {
	client *api.APIClient
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	hostName := d.Get("host_name").(string)
	apiKey := d.Get("api_key").(string)

	config := api.NewConfiguration()

	config.UserAgent = "terraform-provider-mailcow"
	config.Host = hostName
	config.Scheme = "https"
	config.AddDefaultHeader("X-API-Key", apiKey)
	config.AddDefaultHeader("accept", "application/json")
	config.Debug = true

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	apiClient := api.NewAPIClient(config)

	return &APIClient{
		client: apiClient,
	}, diags
}
