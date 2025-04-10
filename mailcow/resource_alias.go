package mailcow

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/l-with/terraform-provider-mailcow/api"
)

func resourceAlias() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAliasCreate,
		ReadContext:   resourceAliasRead,
		UpdateContext: resourceAliasUpdate,
		DeleteContext: resourceAliasDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAliasImport,
		},

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
				Type:             schema.TypeString,
				Description:      `destination address, comma separated. Special values are "ham@localhost", "spam@localhost" and "null@localhost".`,
				Required:         true,
				ValidateDiagFunc: validateGotoSpecialValuesDiag,
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

const (
	gotoHamDestination     = "ham@localhost"
	gotoDiscardDestination = "null@localhost"
	gotoSpamDestination    = "spam@localhost"
)

func validateGotoSpecialValuesDiag(v any, _ cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	gotoValue := v.(string)
	// Special goto values
	specialValues := []string{
		gotoHamDestination,
		gotoDiscardDestination,
		gotoSpamDestination,
	}

	// If the value of goto exactly matches a special value, everything's fine
	if slices.Contains(specialValues, gotoValue) {
		return diags
	}

	// If the value contains a special value but also includes other addresses, it's invalid
	for _, specialValue := range specialValues {
		if strings.Contains(gotoValue, specialValue) {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Invalid value '%s': special value %s is not allowed with other addresses", gotoValue, specialValue),
				Detail:   fmt.Sprintf("The value '%s' cannot contain other addresses along with '%s'. It should only be '%s'.", gotoValue, specialValue, specialValue),
			})
			return diags
		}
	}

	// return if there are just addresses in goto
	return diags
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func resourceAliasImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

func resourceAliasCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*APIClient)

	mailcowCreateRequest := api.NewCreateAliasRequest()

	createRequestSet(mailcowCreateRequest, resourceAlias(), d, nil, nil)

	// if goto is set to one of the special values, set the correspinding goto_ flag and remove the original field
	switch d.Get("goto").(string) {
	case gotoHamDestination:
		mailcowCreateRequest.Set("goto_ham", 1)
		mailcowCreateRequest.Delete("goto")
	case gotoDiscardDestination:
		mailcowCreateRequest.Set("goto_null", 1)
		mailcowCreateRequest.Delete("goto")
	case gotoSpamDestination:
		mailcowCreateRequest.Set("goto_spam", 1)
		mailcowCreateRequest.Delete("goto")
	}

	request := c.client.Api.MailcowCreate(ctx).MailcowCreateRequest(*mailcowCreateRequest)
	response, _, err := c.client.Api.MailcowCreateExecute(request)
	if err != nil {
		return diag.FromErr(err)
	}
	err = checkResponse(response, "resourceAliasCreate", d.Get("address").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := response.GetAliasId()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*id)

	return diags
}

func resourceAliasRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*APIClient)
	id := d.Id()

	request := c.client.Api.MailcowGetAlias(ctx, id)

	alias, err := readRequest(request)
	if err != nil {
		return diag.FromErr(err)
	}

	if alias["id"] == nil {
		return diag.FromErr(errors.New("alias id not found: " + id))
	}

	err = setResourceData(resourceAlias(), d, &alias, nil, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprint(alias["id"].(float64)))

	return diags
}

func resourceAliasUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	mailcowUpdateRequest := api.NewUpdateAliasRequest()

	if d.HasChange("goto") {
		// if goto is set to one of the special values, set the correspinding goto_ flag and remove the original field
		switch d.Get("goto").(string) {
		case gotoHamDestination:
			mailcowUpdateRequest.SetAttr("goto_ham", 1)
			mailcowUpdateRequest.DeleteAttr("goto")
		case gotoDiscardDestination:
			mailcowUpdateRequest.SetAttr("goto_null", 1)
			mailcowUpdateRequest.DeleteAttr("goto")
		case gotoSpamDestination:
			mailcowUpdateRequest.SetAttr("goto_spam", 1)
			mailcowUpdateRequest.DeleteAttr("goto")
		}
	}

	err := mailcowUpdate(ctx, resourceAlias(), d, nil, nil, mailcowUpdateRequest, c)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceAliasRead(ctx, d, m)
}

func resourceAliasDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	mailcowDeleteRequest := api.NewDeleteAliasRequest()
	diags, _ := mailcowDelete(ctx, d, mailcowDeleteRequest, c)
	return diags
}
