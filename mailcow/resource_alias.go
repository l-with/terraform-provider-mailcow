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
	"reflect"
)

func resourceAlias() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAliasCreate,
		ReadContext:   resourceAliasRead,
		UpdateContext: resourceAliasUpdate,
		DeleteContext: resourceAliasDelete,
		Schema: map[string]*schema.Schema{
			"active": {
				Type:        schema.TypeBool,
				Description: "is alias active or not",
				Default:     true,
				Optional:    true,
			},
			"address": {
				Type:        schema.TypeString,
				Description: "alias address, for catchall use \"@domain.tld\"",
				Required:    true,
			},
			"goto": {
				Type:        schema.TypeString,
				Description: "destination address, comma separated",
				Required:    true,
			},
			"goto_ham": {
				Type:        schema.TypeBool,
				Description: "learn as ham",
				Default:     false,
				Optional:    true,
			},
			"goto_null": {
				Type:        schema.TypeBool,
				Description: "silently ignore",
				Default:     false,
				Optional:    true,
			},
			"goto_spam": {
				Type:        schema.TypeBool,
				Description: "learn as spam",
				Default:     false,
				Optional:    true,
			},
			"sogo_visible": {
				Type:        schema.TypeBool,
				Description: "visibility as selectable sender in SOGo",
				Default:     false,
				Optional:    true,
			},
			"private_comment": {
				Type:        schema.TypeString,
				Description: "private comment",
				Optional:    true,
			},
			"public_comment": {
				Type:        schema.TypeString,
				Description: "public comment",
				Optional:    true,
			},
		},
	}
}

func resourceAliasCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*APIClient)

	createAliasRequest := api.NewCreateAliasRequest()
	createAliasRequest.SetActive(d.Get("active").(bool))
	createAliasRequest.SetAddress(d.Get("address").(string))
	createAliasRequest.SetGoto(d.Get("goto").(string))
	createAliasRequest.SetGotoHam(d.Get("goto_ham").(bool))
	createAliasRequest.SetGotoNull(d.Get("goto_null").(bool))
	createAliasRequest.SetGotoSpam(d.Get("goto_spam").(bool))
	createAliasRequest.SetSogoVisible(d.Get("sogo_visible").(bool))
	privateComment, ok := d.GetOk("private_comment")
	if ok {
		createAliasRequest.SetAddress(privateComment.(string))
	}
	publicComment, ok := d.GetOk("public_comment")
	if ok {
		createAliasRequest.SetAddress(publicComment.(string))
	}

	request := c.client.AliasesApi.CreateAlias(ctx).CreateAliasRequest(*createAliasRequest)
	response, _, err := c.client.AliasesApi.CreateAliasExecute(request)
	if err != nil {
		return diag.FromErr(err)
	}

	err = checkResponse(response, resourceAliasCreate, d.Get("address").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	err = checkResponse(response, resourceAliasCreate, d.Get("address").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	err, id := getAliasId(response)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id.(string))

	return diags
}

func getAliasId(response []map[string]interface{}) (error, interface{}) {
	responseMsg := response[len(response)-1]["msg"]
	if reflect.ValueOf(responseMsg).Kind() == reflect.String {
		return errors.New(fmt.Sprint("msg error: ", responseMsg.(string))), nil
	}
	receipt := responseMsg.([]interface{})[0].(string)
	if receipt != "alias_added" {
		return errors.New(fmt.Sprint("msg error: ", receipt)), nil
	}
	return nil, responseMsg.([]interface{})[2].(string)
}

func resourceAliasRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*APIClient)
	id := d.Id()

	request := c.client.AliasesApi.GetAliases(ctx, id)

	response, err := request.Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			diag.FromErr(err)
		}
	}(response.Body)

	alias := make(map[string]interface{}, 0)
	err = json.NewDecoder(response.Body).Decode(&alias)
	if err != nil {
		return diag.FromErr(err)
	}

	if alias["id"] == nil {
		return diag.FromErr(errors.New("alias id not found: " + id))
	}

	for _, argument := range []string{
		"address",
		"goto",
		"private_comment",
		"public_comment",
	} {
		err := d.Set(argument, alias[argument])
		if err != nil {
			return diag.FromErr(err)
		}
	}

	for _, argumentBool := range []string{
		"active",
		"goto_ham",
		"goto_null",
		"goto_spam",
		"sogo_visible",
	} {
		boolValue := false
		if alias[argumentBool] != nil {
			if int(alias[argumentBool].(float64)) >= 1 {
				boolValue = true
			}
		}
		err := d.Set(argumentBool, boolValue)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(fmt.Sprint(alias["id"].(float64)))

	return diags
}

func resourceAliasUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	id := d.Id()

	updateAliasRequest := api.NewUpdateAliasRequest()
	updateAliasRequestAttr := api.NewUpdateAliasRequestAttr()

	if d.HasChange("address") {
		updateAliasRequestAttr.SetActive(d.Get("active").(bool))
	}
	if d.HasChange("address") {
		updateAliasRequestAttr.SetAddress(d.Get("address").(string))
	}
	if d.HasChange("goto") {
		updateAliasRequestAttr.SetAddress(d.Get("goto").(string))
	}
	if d.HasChange("goto_ham") {
		updateAliasRequestAttr.SetGotoHam(d.Get("goto_ham").(bool))
	}
	if d.HasChange("goto_null") {
		updateAliasRequestAttr.SetGotoNull(d.Get("goto_null").(bool))
	}
	if d.HasChange("goto_spam") {
		updateAliasRequestAttr.SetGotoSpam(d.Get("goto_spam").(bool))
	}
	if d.HasChange("sogo_visible") {
		updateAliasRequestAttr.SetSogoVisible(d.Get("sogo_visible").(bool))
	}
	if d.HasChange("private_comment") {
		updateAliasRequestAttr.SetPrivateComment(d.Get("private_comment").(string))
	}
	if d.HasChange("public_comment") {
		updateAliasRequestAttr.SetPrivateComment(d.Get("public_comment").(string))
	}

	items := make([]string, 1)
	items[0] = id

	updateAliasRequest.SetItems(items)
	updateAliasRequest.SetAttr(*updateAliasRequestAttr)
	request := c.client.AliasesApi.UpdateAlias(ctx).UpdateAliasRequest(*updateAliasRequest)
	response, _, err := c.client.AliasesApi.UpdateAliasExecute(request)
	if err != nil {
		return diag.FromErr(err)
	}
	err = checkResponse(response, resourceAliasCreate, d.Get("address").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceAliasRead(ctx, d, m)
}

func resourceAliasDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*APIClient)

	deleteAliasRequest := api.NewDeleteAliasRequest()

	items := make([]string, 1)
	items[0] = d.Id()
	deleteAliasRequest.SetItems(items)

	request := c.client.AliasesApi.DeleteAlias(ctx).DeleteAliasRequest(*deleteAliasRequest)
	response, _, err := c.client.AliasesApi.DeleteAliasExecute(request)
	if err != nil {
		return diag.FromErr(err)
	}
	err = checkResponse(response, resourceAliasDelete, d.Get("address").(string)+" "+d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
