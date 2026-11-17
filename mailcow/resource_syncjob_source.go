package mailcow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/l-with/terraform-provider-mailcow/api"
)

const (
	syncjobSourceScopeAll    = "all"
	syncjobSourceScopeDomain = "domain"
	syncjobSourceScopeUser   = "user"
)

// lowerCaseRegexp matches values without upper case letters (mailcow lower-cases hosts, domains and users)
var lowerCaseRegexp = regexp.MustCompile(`^[^A-Z]*$`)

// syncjobSourceScopeFields only change together: mailcow writes the targets only when the request carries `scope`
var syncjobSourceScopeFields = []string{"scope", "domains", "users"}

func resourceSyncjobSource() *schema.Resource {
	return &schema.Resource{
		Description: "A sync source (mailcow 2026-11 or later) holds the IMAP server that sync jobs copy mail from. " +
			"Sources with OAuth2 (`XOAUTH2`) are not supported.",

		CreateContext: resourceSyncjobSourceCreate,
		ReadContext:   resourceSyncjobSourceRead,
		UpdateContext: resourceSyncjobSourceUpdate,
		DeleteContext: resourceSyncjobSourceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: resourceSyncjobSourceCustomizeDiff,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "name of the sync source, unique per creator",
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
			"description": {
				Type:         schema.TypeString,
				Description:  "description of the sync source",
				Optional:     true,
				Default:      "",
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"host": {
				Type:        schema.TypeString,
				Description: "the IMAP server mails are synced from (--host1), in lower case",
				Required:    true,
				ValidateFunc: validation.All(
					validation.StringIsNotEmpty,
					validation.StringMatch(lowerCaseRegexp, "must be lower case, mailcow stores the host in lower case"),
				),
			},
			"port": {
				Type:         schema.TypeInt,
				Description:  "the port of the IMAP server (--port1)",
				Optional:     true,
				Default:      993,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"enc": {
				Type:         schema.TypeString,
				Description:  "the encryption used to connect to the IMAP server (`SSL`, `TLS` for STARTTLS, `PLAIN`)",
				Optional:     true,
				Default:      "SSL",
				ValidateFunc: validation.StringInSlice([]string{"SSL", "TLS", "PLAIN"}, false),
			},
			"verify_cert": {
				Type:        schema.TypeBool,
				Description: "verify the certificate of the IMAP server",
				Optional:    true,
				Default:     true,
			},
			"auth_type": {
				Type:         schema.TypeString,
				Description:  "the authentication mechanism (`PLAIN`, `LOGIN`, `CRAM-MD5`) (--authmech1)",
				Optional:     true,
				Default:      "PLAIN",
				ValidateFunc: validation.StringInSlice([]string{"PLAIN", "LOGIN", "CRAM-MD5"}, false),
			},
			"scope": {
				Type: schema.TypeString,
				Description: "who may use the source for sync jobs: `all` (everyone, admins only), `domain` (the mailboxes of `domains`) " +
					"or `user` (the mailboxes in `users`, without `users` only the creator and admins)",
				Optional:     true,
				Default:      syncjobSourceScopeUser,
				ValidateFunc: validation.StringInSlice([]string{syncjobSourceScopeAll, syncjobSourceScopeDomain, syncjobSourceScopeUser}, false),
			},
			"domains": {
				Type:        schema.TypeSet,
				Description: "the domains that may use the source, in lower case (only with `scope = \"domain\"`, then required)",
				Optional:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringMatch(lowerCaseRegexp, "must be lower case"),
				},
			},
			"users": {
				Type:        schema.TypeSet,
				Description: "the mailboxes that may use the source, in lower case (only with `scope = \"user\"`)",
				Optional:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringMatch(lowerCaseRegexp, "must be lower case"),
				},
			},
			"active": {
				Type:        schema.TypeBool,
				Description: "is the source active; new sync jobs need an active source",
				Optional:    true,
				Default:     true,
			},
			"created_by": {
				Type:        schema.TypeString,
				Description: "the account that created the source",
				Computed:    true,
			},
		},
	}
}

// resourceSyncjobSourceCustomizeDiff rejects targets that mailcow would drop silently (and that would then show as drift)
func resourceSyncjobSourceCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	scope := d.Get("scope").(string)
	domains := d.Get("domains").(*schema.Set).Len()
	users := d.Get("users").(*schema.Set).Len()
	if scope == syncjobSourceScopeDomain && domains == 0 && d.NewValueKnown("domains") {
		return errors.New(`scope = "domain" needs at least one entry in domains`)
	}
	if scope != syncjobSourceScopeDomain && domains > 0 {
		return errors.New(`domains can only be set with scope = "domain"`)
	}
	if scope != syncjobSourceScopeUser && users > 0 {
		return errors.New(`users can only be set with scope = "user"`)
	}
	return nil
}

// syncjobSourceRequestValue is the value mailcow expects for a schema argument (sets as lists of strings)
func syncjobSourceRequestValue(d *schema.ResourceData, argument string) interface{} {
	value := d.Get(argument)
	if set, ok := value.(*schema.Set); ok {
		list := make([]string, 0, set.Len())
		for _, item := range set.List() {
			list = append(list, item.(string))
		}
		return list
	}
	return value
}

func resourceSyncjobSourceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	name := d.Get("name").(string)

	mailcowCreateRequest := api.NewCreateSyncjobSourceRequest()
	for argument, elem := range resourceSyncjobSource().Schema {
		if elem.Computed && !elem.Optional {
			continue
		}
		mailcowCreateRequest.Set(argument, syncjobSourceRequestValue(d, argument))
	}

	request := c.client.Api.MailcowCreate(ctx).MailcowCreateRequest(*mailcowCreateRequest)
	response, _, err := c.client.Api.MailcowCreateExecute(request)
	if err != nil {
		return diag.FromErr(err)
	}
	err = checkResponse(response, mailcowCreateRequest.ResourceName, name)
	if err != nil {
		return diag.FromErr(err)
	}

	// The response carries no id: the new source is the newest one with this name
	id, err := findSyncjobSourceId(ctx, c, name)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(id))

	return resourceSyncjobSourceRead(ctx, d, m)
}

// findSyncjobSourceId returns the id of the newest visible sync source with the name
func findSyncjobSourceId(ctx context.Context, c *APIClient, name string) (int, error) {
	response, err := c.client.Api.MailcowGetSyncjobSource(ctx, "all").MailcowExecute()
	if err != nil {
		return 0, err
	}
	sources := make([]map[string]interface{}, 0)
	err = json.NewDecoder(response.Body).Decode(&sources)
	if err != nil {
		return 0, err
	}
	id := 0
	for _, source := range sources {
		sourceId, ok := source["id"].(float64)
		if ok && source["name"] == name && int(sourceId) > id {
			id = int(sourceId)
		}
	}
	if id == 0 {
		return 0, fmt.Errorf("syncjob source '%s' not found", name)
	}
	return id, nil
}

// getSyncjobSource returns the sync source, or nil when it does not exist or is not visible
func getSyncjobSource(ctx context.Context, c *APIClient, id string) (map[string]interface{}, error) {
	response, err := c.client.Api.MailcowGetSyncjobSource(ctx, id).MailcowExecute()
	if err != nil {
		return nil, err
	}
	// mailcow returns an empty list for an unknown id, an object otherwise
	var result interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	source, ok := result.(map[string]interface{})
	if !ok || len(source) == 0 {
		return nil, nil
	}
	return source, nil
}

func resourceSyncjobSourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	source, err := getSyncjobSource(ctx, c, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if source == nil {
		log.Printf("[WARN] syncjob source %s not found, removing it from the state", d.Id())
		d.SetId("")
		return nil
	}

	err = setResourceData(resourceSyncjobSource(), d, &source, nil, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceSyncjobSourceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	mailcowUpdateRequest := api.NewUpdateSyncjobSourceRequest()
	for argument, elem := range resourceSyncjobSource().Schema {
		if (elem.Computed && !elem.Optional) || isElementIn(argument, &syncjobSourceScopeFields) {
			continue
		}
		if d.HasChange(argument) {
			mailcowUpdateRequest.SetAttr(argument, syncjobSourceRequestValue(d, argument))
		}
	}
	if d.HasChanges(syncjobSourceScopeFields...) {
		for _, argument := range syncjobSourceScopeFields {
			mailcowUpdateRequest.SetAttr(argument, syncjobSourceRequestValue(d, argument))
		}
	}
	mailcowUpdateRequest.SetItem(d.Id())

	response, err := api.MailcowUpdateExecute(ctx, c.client, mailcowUpdateRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	err = checkResponse(response, mailcowUpdateRequest.ResourceName, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSyncjobSourceRead(ctx, d, m)
}

func resourceSyncjobSourceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	diags, _ := mailcowDelete(ctx, d, api.NewDeleteSyncjobSourceRequest(), c)
	return diags
}
