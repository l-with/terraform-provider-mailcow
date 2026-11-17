package mailcow

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/l-with/terraform-provider-mailcow/api"
)

func resourceSyncjob() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSyncjobCreate,
		ReadContext:   resourceSyncjobRead,
		UpdateContext: resourceSyncjobUpdate,
		DeleteContext: resourceSyncjobDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSyncjobImport,
		},

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type:    resourceSyncjobV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceSyncjobStateUpgradeV0,
			},
		},

		Schema: resourceSyncjobSchema(),
	}
}

func resourceSyncjobSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// mailcow
		"active": {
			Type:        schema.TypeBool,
			Description: "is sync job active or not",
			Default:     true,
			Optional:    true,
		},
		"mins_interval": {
			Type:        schema.TypeInt,
			Description: "the interval in which messages should be synced (minutes)",
			Default:     20,
			Optional:    true,
		},
		"automap": {
			Type:        schema.TypeBool,
			Description: "try to automap folders (\"Sent items\", \"Sent\" => \"Sent\" etc.) (--automap)",
			Default:     true,
			Optional:    true,
		},

		// imapsync
		"custom_params": {
			Type:        schema.TypeList,
			Description: "custom imapsync parameters; mailcow only accepts the options of its allowlist",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"option": {
						Type:        schema.TypeString,
						Description: "the imapsync option, lower case and without leading dashes (e.g. `nofoldersizes`)",
						Required:    true,
						ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-z][a-z0-9]*$`),
							"must be lower case letters and digits, without leading dashes"),
					},
					"value": {
						Type:        schema.TypeString,
						Description: "the value of the option, empty for a flag",
						Optional:    true,
						Default:     "",
					},
				},
			},
		},
		"delete1": {
			Type:        schema.TypeBool,
			Description: "delete (mail) from source when completed (--delete1)",
			Default:     false,
			Optional:    true,
		},
		"delete2": {
			Type:        schema.TypeBool,
			Description: "delete messages on destination that are not on source (--delete2)",
			Default:     false,
			Optional:    true,
		},
		"delete2duplicates": {
			Type:        schema.TypeBool,
			Description: "delete duplicates on destination (--delete2duplicates)",
			Default:     true,
			Optional:    true,
		},
		"exclude": {
			Type:        schema.TypeString,
			Description: "exclude objects (regex) (--exclude)",
			Default:     "",
			Optional:    true,
		},
		"maxage": {
			Type:        schema.TypeInt,
			Description: "only sync messages up to this age in days (--maxage)",
			Default:     0,
			Optional:    true,
		},
		"maxbytespersecond": {
			Type:        schema.TypeString,
			Description: "max speed transfer limit for the sync (--maxbytespersecond)",
			Default:     0,
			Optional:    true,
		},
		"skipcrossduplicates": {
			Type:        schema.TypeBool,
			Description: "skip duplicate messages across folders (first come, first serve) (--skipcrossduplicates)",
			Default:     false,
			Optional:    true,
		},
		"subscribeall": {
			Type:        schema.TypeBool,
			Description: "subscribe all folders (--subscribeall)",
			Default:     true,
			Optional:    true,
		},
		"subfolder2": {
			Type:        schema.TypeString,
			Description: "sync into subfolder on destination (empty = do not use subfolder) (--subfolder2)",
			Default:     "",
			Optional:    true,
		},
		"timeout2": {
			Type:        schema.TypeInt,
			Description: "timeout for connection to local host (--timeout2)",
			Default:     600,
			Optional:    true,
		},

		// imapsync source (host1)
		"source_id": {
			Type:        schema.TypeInt,
			Description: "the id of the sync source (e.g. `mailcow_syncjob_source.source.id`) with the IMAP server mails are synced from",
			Required:    true,
		},
		"password1": {
			Type:        schema.TypeString,
			Description: "the password of the mailbox on the host (--password1)",
			Required:    true,
			Sensitive:   true,
		},
		"timeout1": {
			Type:        schema.TypeInt,
			Description: "timeout for connection to remote host (--timeout1)",
			Default:     600,
			Optional:    true,
		},
		"user1": {
			Type:        schema.TypeString,
			Description: "user to login on remote host (--user1)",
			Required:    true,
		},
		"username": { // user2 on get
			Type:        schema.TypeString,
			Description: "user to login on local host (--user2)",
			Required:    true,
			ForceNew:    true,
		},
	}
}

func resourceSyncjobImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

// syncjobCustomParamsExclude are not sent by the generic create and update helpers, custom_params goes as a list of pairs
var syncjobCustomParamsExclude = []string{"custom_params"}

// syncjobCustomParamsRequest converts custom_params to the {o: option, v: value} pairs mailcow expects
func syncjobCustomParamsRequest(d *schema.ResourceData) []map[string]string {
	pairs := make([]map[string]string, 0)
	for _, item := range d.Get("custom_params").([]interface{}) {
		if item == nil {
			continue
		}
		param := item.(map[string]interface{})
		pairs = append(pairs, map[string]string{"o": param["option"].(string), "v": param["value"].(string)})
	}
	return pairs
}

// syncjobCustomParamsState converts the custom_params JSON of mailcow ({o, v} pairs) to the state
func syncjobCustomParamsState(value interface{}) []map[string]interface{} {
	params := make([]map[string]interface{}, 0)
	text, ok := value.(string)
	if !ok || text == "" {
		return params
	}
	pairs := make([]map[string]interface{}, 0)
	if err := json.Unmarshal([]byte(text), &pairs); err != nil {
		log.Print("[WARN] syncjob custom_params is not a list of pairs: ", text)
		return params
	}
	for _, pair := range pairs {
		params = append(params, map[string]interface{}{"option": fmt.Sprint(pair["o"]), "value": fmt.Sprint(pair["v"])})
	}
	return params
}

func resourceSyncjobCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*APIClient)

	mailcowCreateRequest := api.NewCreateSyncjobRequest()

	createRequestSet(mailcowCreateRequest, resourceSyncjob(), d, &syncjobCustomParamsExclude, nil)
	mailcowCreateRequest.Set("custom_params", syncjobCustomParamsRequest(d))

	username := d.Get("username").(string)
	user1 := d.Get("user1").(string)
	sourceId := d.Get("source_id").(int)

	request := c.client.Api.MailcowCreate(ctx).MailcowCreateRequest(*mailcowCreateRequest)
	response, _, err := c.client.Api.MailcowCreateExecute(request)
	if err != nil {
		return diag.FromErr(err)
	}
	err = checkResponse(response, "resourceSyncjobCreate", username+"=>"+user1)
	if err != nil {
		return diag.FromErr(err)
	}

	// (source, local mailbox, remote user) is unique
	syncJob, err := getSyncJob(ctx, c, username, func(job map[string]interface{}) bool {
		return job["user1"] == user1 && fmt.Sprint(job["source_id"]) == fmt.Sprint(sourceId)
	})
	if syncJob == nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprint(syncJob["id"].(float64)))

	return diags
}

func resourceSyncjobRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var err error

	c := m.(*APIClient)
	id := d.Id()
	emailAddress := d.Get("username").(string)

	syncJob, err := getSyncJob(ctx, c, emailAddress, func(job map[string]interface{}) bool {
		return fmt.Sprint(job["id"]) == id
	})
	if syncJob == nil {
		return diag.FromErr(err)
	}

	syncJob["username"] = syncJob["user2"]
	for _, argument := range []string{
		"mins_interval",
		"exclude",
		"maxage",
		"maxbytespersecond",
		"subfolder2",
		"timeout2",
		"source_id",
		"timeout1",
		"user1",
		"username",
	} {
		err = d.Set(argument, syncJob[argument])
		if err != nil {
			return diag.FromErr(err)
		}
		log.Print("[TRACE] resourceSyncjobRead mailbox[", argument, "]: ", syncJob[argument])
	}

	err = d.Set("custom_params", syncjobCustomParamsState(syncJob["custom_params"]))
	if err != nil {
		return diag.FromErr(err)
	}

	for _, argumentBool := range []string{
		"active",
		"automap",
		"delete1",
		"delete2",
		"delete2duplicates",
		"skipcrossduplicates",
		"subscribeall",
	} {
		boolValue := false
		if syncJob[argumentBool] != nil {
			if int(syncJob[argumentBool].(float64)) >= 1 {
				boolValue = true
			}
		}
		err = d.Set(argumentBool, boolValue)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(id)

	return diags
}

// getSyncJob returns the first sync job of the mailbox that matches
func getSyncJob(ctx context.Context, c *APIClient, emailAddress string, match func(map[string]interface{}) bool) (map[string]interface{}, error) {
	request := c.client.Api.MailcowGetSyncjob(ctx, emailAddress)
	log.Print("[TRACE] getSyncJob emailAddress: ", emailAddress)

	response, err := request.MailcowExecute()
	if err != nil {
		return nil, err
	}

	log.Print("[TRACE] getSyncJob response.Body: ", response.Body)
	syncJobs := make([]map[string]interface{}, 0)
	err = json.NewDecoder(response.Body).Decode(&syncJobs)
	if err != nil {
		return nil, err
	}

	for _, syncJob := range syncJobs {
		if syncJob["user2"] == emailAddress && match(syncJob) {
			return syncJob, nil
		}
	}
	return nil, fmt.Errorf("syncjob of %s not found", emailAddress)
}

func resourceSyncjobUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	mailcowUpdateRequest := api.NewUpdateSyncjobRequest()
	if d.HasChange("custom_params") {
		mailcowUpdateRequest.SetAttr("custom_params", syncjobCustomParamsRequest(d))
	}

	err := mailcowUpdate(ctx, resourceSyncjob(), d, &syncjobCustomParamsExclude, nil, mailcowUpdateRequest, c)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSyncjobRead(ctx, d, m)
}

func resourceSyncjobDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	mailcowDeleteRequest := api.NewDeleteSyncjobRequest()
	diags, _ := mailcowDelete(ctx, d, mailcowDeleteRequest, c)
	return diags
}

// resourceSyncjobV0 is the schema before mailcow 2026-11: the IMAP server in host1/port1/enc1, custom_params as text
func resourceSyncjobV0() *schema.Resource {
	s := resourceSyncjobSchema()
	delete(s, "source_id")
	s["custom_params"] = &schema.Schema{Type: schema.TypeString, Optional: true}
	s["host1"] = &schema.Schema{Type: schema.TypeString, Required: true}
	s["port1"] = &schema.Schema{Type: schema.TypeInt, Optional: true}
	s["enc1"] = &schema.Schema{Type: schema.TypeString, Optional: true}
	return &schema.Resource{Schema: s}
}

// resourceSyncjobStateUpgradeV0 drops host1/port1/enc1 (now in the sync source) and the text custom_params;
// source_id and custom_params are read from mailcow on the next refresh
func resourceSyncjobStateUpgradeV0(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	if rawState == nil {
		return rawState, nil
	}
	delete(rawState, "host1")
	delete(rawState, "port1")
	delete(rawState, "enc1")
	rawState["custom_params"] = []interface{}{}
	return rawState, nil
}
