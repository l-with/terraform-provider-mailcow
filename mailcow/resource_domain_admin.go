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

func resourceDomainAdmin() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainAdminCreate,
		ReadContext:   resourceDomainAdminRead,
		UpdateContext: resourceDomainAdminUpdate,
		DeleteContext: resourceDomainAdminDelete,
		Schema: map[string]*schema.Schema{
			"active": {
				Type:        schema.TypeBool,
				Description: "is domain admin active or not",
				Default:     true,
				Optional:    true,
			},
			"username": {
				Type:        schema.TypeString,
				Description: "username of the domain admin",
				Required:    true,
				ForceNew:    true,
			},
			"domains": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "domain names",
				Required:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "domain admin password",
				Required:    true,
				Sensitive:   true,
			},
		},
	}
}

func resourceDomainAdminCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*APIClient)

	mailcowCreateRequest := api.NewCreateDomainAdminRequest()

	username := d.Get("username").(string)
	mailcowCreateRequest.Set("password2", d.Get("password"))

	err := mailcowCreate(ctx, resourceDomainAdmin(), d, username, nil, nil, mailcowCreateRequest, c)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(username)

	return diags
}

func resourceDomainAdminRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var err error

	c := m.(*APIClient)
	id := d.Id()
	username := d.Get("username").(string)

	syncJob, err := getSyncJob(ctx, c, username, "id", id)
	if syncJob == nil {
		return diag.FromErr(err)
	}

	for _, argument := range []string{
		"domains",
	} {
		err = d.Set(argument, syncJob[argument])
		if err != nil {
			return diag.FromErr(err)
		}
		log.Print("[TRACE] resourceSyncjobRead mailbox[", argument, "]: ", syncJob[argument])
	}

	for _, argumentBool := range []string{
		"active",
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

func getDomainAdmin(ctx context.Context, c *APIClient, name string, attributeKey string, attributeValue string) (map[string]interface{}, error) {
	request := c.client.Api.MailcowGetDomainAdmin(ctx, name)
	log.Print("[TRACE] getDomainAdmin name: ", name)

	response, err := request.MailcowExecute()
	if err != nil {
		return nil, err
	}

	log.Print("[TRACE] getDomainAdmin response.Body: ", response.Body)
	domainAdmins := make([]map[string]interface{}, 0)
	err = json.NewDecoder(response.Body).Decode(&domainAdmins)
	if err != nil {
		return nil, err
	}

	var syncJob map[string]interface{} = nil
	for _, currentDomainAdmin := range domainAdmins {
		currentUsername := currentDomainAdmin["username"].(string)
		if currentUsername == name && attributeValue == fmt.Sprint(currentDomainAdmin[attributeKey]) {
			syncJob = currentDomainAdmin
			break
		}
	}
	if syncJob == nil {
		return nil, errors.New(fmt.Sprintf("domain-admin username=%s and %s=%s not found", name, attributeKey, attributeValue))
	}
	return syncJob, nil
}

func resourceDomainAdminUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	mailcowUpdateRequest := api.NewUpdateDomainAdminRequest()

	exclude := []string{
		"password",
	}
	err := mailcowUpdate(ctx, resourceDomainAdmin(), d, &exclude, nil, mailcowUpdateRequest, c)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDomainAdminRead(ctx, d, m)
}

func resourceDomainAdminDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	mailcowDeleteRequest := api.NewDeleteDomainAdminRequest()
	diags, _ := mailcowDelete(ctx, d, mailcowDeleteRequest, c)
	return diags
}
