package mailcow

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/l-with/terraform-provider-mailcow/api"
)

type MailcowRelayhost struct {
	Id              int
	Hostname        string
	Username        string
	Password        string
	PasswordShort   string
	UsedByDomains   string
	UsedByMailboxes string
}

func resourceRelayhost() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a relayhost / sender dependent transport configuration in mailcow. This can be used in order to use third party SMTP-Relays with better reputation such as Sweego, Lettermint or AWS SES.",

		CreateContext: resourceRelayhostCreate,
		ReadContext:   resourceRelayhostRead,
		UpdateContext: resourceRelayhostUpdate,
		DeleteContext: resourceRelayhostDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:        schema.TypeString,
				Description: "The hostname to connect to, including port. Must be of the format {hostname}:{port}.",
				Required:    true,
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					val, ok := value.(string)
					if !ok {
						return diag.FromErr(fmt.Errorf("expected a string, got %T", value))
					}

					if !strings.Contains(val, ".") || !strings.Contains(val, ":") {
						return diag.FromErr(fmt.Errorf(
							"the hostname of a relayhost must be a hostname or an IP-Address followed by a port, separated by a colon (e.g. `foo.com:123`)",
						))
					}

					return nil
				},
			},
			"username": {
				Type:        schema.TypeString,
				Description: "The username to use for the connection",
				Required:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "The password to use for the connection",
				Required:    true,
				Sensitive:   true,
			},
		},
	}
}

func resourceRelayhostCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*APIClient)
	hostname := d.Get("hostname").(string)
	username := d.Get("username").(string)

	existingId, _ := getRelayhostId(ctx, c, hostname, username)
	if existingId != -1 {
		return diag.FromErr(fmt.Errorf("relayhost with hostname=%s and username=%s already exists", hostname, username))
	}

	mailcowCreateRequest := api.NewCreateRelayhostRequest()
	createRequestSet(mailcowCreateRequest, resourceRelayhost(), d, nil, nil)

	request := c.client.Api.MailcowCreate(ctx).MailcowCreateRequest(*mailcowCreateRequest)
	response, _, err := c.client.Api.MailcowCreateExecute(request)
	if err != nil {
		return diag.FromErr(err)
	}

	err = checkResponse(response, "resourceRelayhostCreate", hostname)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := getRelayhostId(ctx, c, hostname, username)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", id))

	return diags
}

func resourceRelayhostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	id := d.Id()

	relayhost, err := getRelayhost(ctx, c, id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = mapRelayhost(*relayhost, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func resourceRelayhostUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	updateRequest := api.NewUpdateRelayhostRequest()
	err := mailcowUpdate(ctx, resourceRelayhost(), d, nil, nil, updateRequest, c)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRelayhostRead(ctx, d, m)
}

func resourceRelayhostDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	deleteRequest := api.NewDeleteRelayhostRequest()
	diags, _ := mailcowDelete(ctx, d, deleteRequest, c)
	return diags
}

func getRelayhost(ctx context.Context, c *APIClient, id string) (*MailcowRelayhost, error) {
	request := c.client.Api.MailcowGetRelayhost(ctx, id)
	log.Printf("[DEBUG] getRelayhost id: %d", id)

	response, err := request.MailcowExecute()
	if err != nil {
		return nil, err
	}

	relayhost := &MailcowRelayhost{}
	err = json.NewDecoder(response.Body).Decode(relayhost)
	if err != nil {
		return nil, err
	}

	// API Always returns 200, but an empty JSON object on error, so we need
	// to detect that
	if relayhost.Hostname == "" {
		return nil, fmt.Errorf("Relayhost id=%d not found", id)
	}

	return relayhost, nil
}

func getRelayhostId(ctx context.Context, c *APIClient, hostname string, username string) (int, error) {
	request := c.client.Api.MailcowGetRelayhosts(ctx)
	log.Printf("[TRACE] getRelayhostId hostname: %s username: %s", hostname, username)

	response, err := request.MailcowExecute()
	if err != nil {
		return -1, err
	}

	log.Printf("[TRACE] getRelayhost response.Body: ", response.Body)
	relayhosts := make([]MailcowRelayhost, 0)
	err = json.NewDecoder(response.Body).Decode(&relayhosts)
	if err != nil {
		return -1, err
	}

	for _, relayhost := range relayhosts {
		if relayhost.Hostname == hostname && relayhost.Username == username {
			return relayhost.Id, nil
		}
	}

	return -1, fmt.Errorf("relayhost hostname=%s username=%s not found", hostname, username)
}

func mapRelayhost(relayhost MailcowRelayhost, d *schema.ResourceData) error {
	properties := map[string]interface{}{
		"hostname": relayhost.Hostname,
		"username": relayhost.Username,
		"password": relayhost.Password,
	}

	d.SetId(fmt.Sprintf("%d", relayhost.Id))
	for key, value := range properties {
		err := d.Set(key, value)
		if err != nil {
			return fmt.Errorf("Cannot map relayhost %d into properties: %s", relayhost.Id, err)
		}
	}

	return nil
}
