package mailcow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "github.com/l-with/mailcow-go"
	"io"
	"log"
	"strconv"
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
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*APIClient)

	createSyncjobRequest := api.NewCreateSyncJobRequest()
	createSyncjobRequest.SetActive(d.Get("active").(bool))
	createSyncjobRequest.SetMinsInterval(float32(d.Get("mins_interval").(int)))
	createSyncjobRequest.SetAutomap(d.Get("automap").(bool))
	createSyncjobRequest.SetCustomParams(d.Get("custom_params").(string))
	createSyncjobRequest.SetDelete1(d.Get("delete1").(bool))
	createSyncjobRequest.SetDelete2(d.Get("delete2").(bool))
	createSyncjobRequest.SetDelete2duplicates(d.Get("delete2duplicates").(bool))
	createSyncjobRequest.SetExclude(d.Get("exclude").(string))
	createSyncjobRequest.SetMaxage(float32(d.Get("maxage").(int)))
	maxBytesPerSecondNumber, err := strconv.Atoi(d.Get("maxbytespersecond").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	createSyncjobRequest.SetMaxbytespersecond(float32(maxBytesPerSecondNumber))
	createSyncjobRequest.SetSkipcrossduplicates(d.Get("skipcrossduplicates").(bool))
	createSyncjobRequest.SetSubscribeall(d.Get("subscribeall").(bool))
	createSyncjobRequest.SetSubfolder2(d.Get("subfolder2").(string))
	createSyncjobRequest.SetTimeout2(float32(d.Get("timeout2").(int)))
	createSyncjobRequest.SetEnc1(d.Get("enc1").(string))
	createSyncjobRequest.SetHost1(d.Get("host1").(string))
	createSyncjobRequest.SetPassword1(d.Get("password1").(string))
	createSyncjobRequest.SetPort1(float32(d.Get("port1").(int)))
	createSyncjobRequest.SetTimeout1(float32(d.Get("timeout1").(int)))
	createSyncjobRequest.SetUser1(d.Get("user1").(string))
	createSyncjobRequest.SetUsername(d.Get("username").(string))

	username := d.Get("username").(string)
	user1 := d.Get("user1").(string)
	request := c.client.SyncJobsApi.CreateSyncJob(ctx).CreateSyncJobRequest(*createSyncjobRequest)
	response, _, err := c.client.SyncJobsApi.CreateSyncJobExecute(request)
	if err != nil {
		return diag.FromErr(err)
	}
	err = checkResponse(response, resourceSyncjobCreate, username+"=>"+user1)
	if err != nil {
		return diag.FromErr(err)
	}

	diags, syncJob := getSyncJob(ctx, c, username, "user1", d.Get("user1").(string))
	if syncJob == nil {
		return diags
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

	diags, syncJob := getSyncJob(ctx, c, emailAddress, "id", id)
	if syncJob == nil {
		return diags
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

func getSyncJob(ctx context.Context, c *APIClient, emailAddress string, attributeKey string, attributeValue string) (diag.Diagnostics, map[string]interface{}) {
	request := c.client.SyncJobsApi.GetSyncJobs(ctx, emailAddress)

	response, err := request.Execute()
	if err != nil {
		return diag.FromErr(err), nil
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			diag.FromErr(err)
		}
	}(response.Body)

	syncJobs := make([]map[string]interface{}, 0)
	err = json.NewDecoder(response.Body).Decode(&syncJobs)
	if err != nil {
		return diag.FromErr(err), nil
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
		return diag.FromErr(errors.New(fmt.Sprintf("syncjob user2=%s and %s=%s not found", emailAddress, attributeKey, attributeValue))), nil
	}
	return nil, syncJob
}

func resourceSyncjobUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	id := d.Id()

	updateSyncjobRequest := api.NewUpdateSyncJobRequest()
	updateSyncjobRequestAttr := api.NewUpdateSyncJobRequestAttr()

	if d.HasChange("active") {
		updateSyncjobRequestAttr.SetActive(d.Get("active").(bool))
	}
	if d.HasChange("mins_interval") {
		updateSyncjobRequestAttr.SetMinsInterval(float32(d.Get("mins_interval").(int)))
	}
	if d.HasChange("automap") {
		updateSyncjobRequestAttr.SetAutomap(d.Get("automap").(bool))
	}
	if d.HasChange("custom_params") {
		updateSyncjobRequestAttr.SetCustomParams(d.Get("custom_params").(string))
	}
	if d.HasChange("delete1") {
		updateSyncjobRequestAttr.SetDelete1(d.Get("delete1").(bool))
	}
	if d.HasChange("delete2") {
		updateSyncjobRequestAttr.SetDelete2(d.Get("delete2").(bool))
	}
	if d.HasChange("delete2duplicates") {
		updateSyncjobRequestAttr.SetDelete2duplicates(d.Get("delete2duplicates").(bool))
	}
	if d.HasChange("exclude") {
		updateSyncjobRequestAttr.SetExclude(d.Get("exclude").(string))
	}
	if d.HasChange("maxage") {
		updateSyncjobRequestAttr.SetMaxage(float32(d.Get("maxage").(int)))
	}
	if d.HasChange("maxbytespersecond") {
		maxBytesPerSecondNumber, err := strconv.Atoi(d.Get("maxbytespersecond").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		updateSyncjobRequestAttr.SetMaxbytespersecond(float32(maxBytesPerSecondNumber))
	}
	if d.HasChange("skipcrossduplicates") {
		updateSyncjobRequestAttr.SetSkipcrossduplicates(d.Get("skipcrossduplicates").(bool))
	}
	if d.HasChange("subscribeall") {
		updateSyncjobRequestAttr.SetSubscribeall(d.Get("subscribeall").(bool))
	}
	if d.HasChange("subfolder2") {
		updateSyncjobRequestAttr.SetSubfolder2(d.Get("subfolder2").(string))
	}
	if d.HasChange("timeout2") {
		updateSyncjobRequestAttr.SetTimeout2(float32(d.Get("timeout2").(int)))
	}
	if d.HasChange("enc1") {
		updateSyncjobRequestAttr.SetEnc1(d.Get("enc1").(string))
	}
	if d.HasChange("host1") {
		updateSyncjobRequestAttr.SetHost1(d.Get("host1").(string))
	}
	if d.HasChange("password1") {
		updateSyncjobRequestAttr.SetPassword1(d.Get("password1").(string))
	}
	if d.HasChange("port1") {
		updateSyncjobRequestAttr.SetPort1(float32(d.Get("port1").(int)))
	}
	if d.HasChange("timeout1") {
		updateSyncjobRequestAttr.SetTimeout1(float32(d.Get("timeout1").(int)))
	}
	if d.HasChange("user1") {
		updateSyncjobRequestAttr.SetUser1(d.Get("user1").(string))
	}

	items := make([]string, 1)
	items[0] = id

	updateSyncjobRequest.SetItems(items)
	updateSyncjobRequest.SetAttr(*updateSyncjobRequestAttr)
	request := c.client.SyncJobsApi.UpdateSyncJob(ctx).UpdateSyncJobRequest(*updateSyncjobRequest)
	_, _, err := c.client.SyncJobsApi.UpdateSyncJobExecute(request)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSyncjobRead(ctx, d, m)
}

func resourceSyncjobDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*APIClient)

	deleteSyncjobRequest := api.NewDeleteSyncJobRequest()

	items := make([]string, 1)
	items[0] = d.Id()
	deleteSyncjobRequest.SetItems(items)

	request := c.client.SyncJobsApi.DeleteSyncJob(ctx).DeleteSyncJobRequest(*deleteSyncjobRequest)
	response, _, err := c.client.SyncJobsApi.DeleteSyncJobExecute(request)
	if err != nil {
		return diag.FromErr(err)
	}
	if response[len(response)-1]["type"].(string) != "success" {
		return diag.FromErr(errors.New(response[0]["type"].(string)))
	}

	d.SetId("")

	return diags
}
