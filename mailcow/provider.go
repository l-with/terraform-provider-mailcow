package mailcow

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	mailcow "github.com/l-with/mailcow-go"
)

var stderr = os.Stderr

func New() tfsdk.Provider {
	return &provider{}
}

type provider struct {
	configured bool
	client     *mailcow.APIClient
}

// GetSchema
func (p *provider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"url": {
				Type:        types.StringType,
				Required:    true,
				Description: "The mailcow API endpoint, can optionally be passed as `MAILCOW_API_URL` environmental variable",
			},
			"token": {
				Type:        types.StringType,
				Required:    true,
				Sensitive:   true,
				Description: "The mailcow API key, can optionally be passed as `MAILCOW_API_KEY` environmental variable",
			},
		},
	}, nil
}

// Provider schema struct
type providerData struct {
	Url   types.String `tfsdk:"url"`
	Token types.String `tfsdk:"token"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	// Retrieve provider data from configuration
	var providerConfig providerData
	diags := req.Config.Get(ctx, &providerConfig)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// User must provide a user to the provider
	var apiUrl string
	if providerConfig.Url.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as url",
		)
		return
	}

	if providerConfig.Url.Null {
		apiUrl = os.Getenv("MAILCOW_API_URL")
	} else {
		apiUrl = providerConfig.Url.Value
	}

	if apiUrl == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find url",
			"Url cannot be an empty string",
		)
		return
	}

	// User must provide a password to the provider
	var token string
	if providerConfig.Token.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Cannot use unknown value as token",
		)
		return
	}

	if providerConfig.Token.Null {
		token = os.Getenv("MAILCOW_API_KEY")
	} else {
		token = providerConfig.Token.Value
	}

	if token == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find token",
			"token cannot be an empty string",
		)
		return
	}

	mcUrl, err := url.Parse(apiUrl)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	config := api.NewConfiguration()
	config.Debug = true
	config.UserAgent = fmt.Sprintf("mailcow-terraform@%s", version)
	config.Host = mcUrl.Host
	config.Scheme = mcUrl.Scheme

	// Create a new  client and set it to the provider client
	c, err := api.NewClient(&host, &username, &password)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Unable to create hashicups client:\n\n"+err.Error(),
		)
		return
	}

	p.client = c
	p.configured = true
}

// func providerConfigure(version string, testing bool) schema.ConfigureContextFunc {
// 	return func(c context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
// 		apiURL := d.Get("url").(string)
// 		token := d.Get("token").(string)
// 		insecure := d.Get("insecure").(bool)

// 		// Warning or errors can be collected in a slice type
// 		var diags diag.Diagnostics

// 		akURL, err := url.Parse(apiURL)
// 		if err != nil {
// 			return nil, diag.FromErr(err)
// 		}

// 		config := api.NewConfiguration()
// 		config.Debug = true
// 		config.UserAgent = fmt.Sprintf("authentik-terraform@%s", version)
// 		config.Host = akURL.Host
// 		config.Scheme = akURL.Scheme
// 		if testing {
// 			config.HTTPClient = &http.Client{
// 				Transport: NewTestingTransport(GetTLSTransport(insecure)),
// 			}
// 		} else {
// 			config.HTTPClient = &http.Client{
// 				Transport: GetTLSTransport(insecure),
// 			}
// 		}

// 		config.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", token))
// 		apiClient := api.NewAPIClient(config)

// 		rootConfig, _, err := apiClient.RootApi.RootConfigRetrieve(context.Background()).Execute()
// 		if err == nil && rootConfig.ErrorReporting.Enabled {
// 			dsn := "https://7b485fd979bf48c1acbe38ffe382a541@sentry.beryju.org/14"
// 			if envDsn, found := os.LookupEnv("SENTRY_DSN"); found {
// 				dsn = envDsn
// 			}
// 			err := sentry.Init(sentry.ClientOptions{
// 				Dsn:              dsn,
// 				Environment:      rootConfig.ErrorReporting.Environment,
// 				TracesSampleRate: float64(rootConfig.ErrorReporting.TracesSampleRate),
// 				Release:          fmt.Sprintf("authentik-terraform-provider@%s", version),
// 			})
// 			if err != nil {
// 				fmt.Printf("Error during sentry init: %v\n", err)
// 			}
// 			config.HTTPClient.Transport = NewTracingTransport(context.Background(), config.HTTPClient.Transport)
// 			apiClient = api.NewAPIClient(config)
// 		}

// 		return &APIClient{
// 			client: apiClient,
// 		}, diags
// 	}
// }

// GetResources - Defines provider resources
func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{}, nil
}

func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{}, nil
}
