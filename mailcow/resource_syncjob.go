package mailcow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/l-with/terraform-provider-mailcow/api"
	"log"
)

func resourceSyncjob() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSyncjobCreate,
		ReadContext:   resourceSyncjobRead,
		UpdateContext: resourceSyncjobUpdate,
		DeleteContext: resourceSyncjobDelete,
		Schema: map[string]*schema.Schema{
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
				Type:        schema.TypeString,
				Description: "custom parameters",
				Default:     "",
				Optional:    true,
			},
			"delete1": {
				Type:        schema.TypeBool,
				Description: "delete from source when completed (--delete1)",
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

			// imapsync target (host1)
			"enc1": {
				Type:        schema.TypeString,
				Description: "the encryption method used to connect to the target mailserver (SSL,TLS,PLAIN)",
				Default:     "SSL",
				Optional:    true,
			},
			"host1": {
				Type:        schema.TypeString,
				Description: "the smtp server where mails should be synced from (--host1)",
				Required:    true,
			},
			"password1": {
				Type:        schema.TypeString,
				Description: "the password of the mailbox on the host (--password1)",
				Required:    true,
				Sensitive:   true,
			},
			"port1": {
				Type:        schema.TypeInt, // in openapi spec string
				Description: "the smtp port of the target mail server (--port1)",
				Default:     143,
				Optional:    true,
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
		},
	}
}

func resourceSyncjobCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*APIClient)

	mailcowCreateRequest := api.NewCreateSyncjobRequest()

	log.Print("[TRACE] resourceSyncjobCreate delete1: ", d.Get("delete1"))
	createRequestSet(mailcowCreateRequest, resourceSyncjob(), d, nil, nil)

	username := d.Get("username").(string)
	user1 := d.Get("user1").(string)

	request := c.client.Api.MailcowCreate(ctx).MailcowCreateRequest(*mailcowCreateRequest)
	response, _, err := c.client.Api.MailcowCreateExecute(request)
	if err != nil {
		return diag.FromErr(err)
	}
	err = checkResponse(response, "resourceSyncjobCreate", username+"=>"+user1)
	if err != nil {
		return diag.FromErr(err)
	}

	syncJob, err := getSyncJob(ctx, c, username, "user1", d.Get("user1").(string))
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

	syncJob, err := getSyncJob(ctx, c, emailAddress, "id", id)
	if syncJob == nil {
		return diag.FromErr(err)
	}

	syncJob["username"] = syncJob["user2"]
	for _, argument := range []string{
		"mins_interval",
		"custom_params",
		"exclude",
		"maxage",
		"maxbytespersecond",
		"subfolder2",
		"timeout2",
		"enc1",
		"host1",
		"port1",
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

func getSyncJob(ctx context.Context, c *APIClient, emailAddress string, attributeKey string, attributeValue string) (map[string]interface{}, error) {
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

	var syncJob map[string]interface{} = nil
	for _, currentSyncJob := range syncJobs {
		currentUsername := currentSyncJob["user2"].(string)
		if currentUsername == emailAddress && attributeValue == fmt.Sprint(currentSyncJob[attributeKey]) {
			syncJob = currentSyncJob
			break
		}
	}
	if syncJob == nil {
		return nil, errors.New(fmt.Sprintf("syncjob user2=%s and %s=%s not found", emailAddress, attributeKey, attributeValue))
	}
	return syncJob, nil
}

func resourceSyncjobUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	mailcowUpdateRequest := api.NewUpdateSyncjobRequest()

	err := mailcowUpdate(ctx, resourceSyncjob(), d, nil, nil, mailcowUpdateRequest, c)
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
