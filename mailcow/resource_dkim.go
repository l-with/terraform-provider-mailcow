package mailcow

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "github.com/l-with/mailcow-go"
	"io"
	"reflect"
	"strconv"
)

func resourceDkim() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDkimCreate,
		ReadContext:   resourceDkimRead,
		DeleteContext: resourceDkimDelete,
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"pubkey": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"length": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"dkim_txt": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dkim_selector": {
				Type:     schema.TypeString,
				Default:  "dkim",
				Optional: true,
				ForceNew: true,
			},
			//"privkey": {
			//	Type:      schema.TypeString,
			//	Computed:  true,
			//	Sensitive: true,
			//},
		},
	}
}

func resourceDkimCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	// var diags diag.Diagnostics

	c := m.(*APIClient)

	createDkimRequest := api.NewGenerateDKIMKeyRequest()
	createDkimRequest.SetDomains(d.Get("domain").(string))
	createDkimRequest.SetKeySize(float32(d.Get("length").(int)))
	createDkimRequest.SetDkimSelector(d.Get("dkim_selector").(string))

	request := c.client.DKIMApi.GenerateDKIMKey(ctx).GenerateDKIMKeyRequest(*createDkimRequest)
	response, _, err := c.client.DKIMApi.GenerateDKIMKeyExecute(request)
	if err != nil {
		return diag.FromErr(err)
	}
	err = checkResponse(response, resourceDkimCreate, d.Get("domain").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("domain").(string))

	return resourceDkimRead(ctx, d, m)
}

func resourceDkimRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*APIClient)
	domain := d.Id()

	request := c.client.DKIMApi.GetDKIMKey(ctx, domain)

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

	dkim := make(map[string]interface{}, 0)
	err = json.NewDecoder(response.Body).Decode(&dkim)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, argument := range []string{
		"pubkey",
		"dkim_txt",
		"dkim_selector",
		//"privkey",
	} {
		dkimArgument := dkim[argument]
		if dkimArgument != nil {
			err = d.Set(argument, dkimArgument)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	for _, argumentNumber := range []string{
		"length",
	} {
		dkimArgumentNumber := dkim[argumentNumber]
		if dkimArgumentNumber != nil {
			value := reflect.ValueOf(dkim[argumentNumber])
			intValue, err := strconv.Atoi(fmt.Sprint(value))
			err = d.Set(argumentNumber, intValue)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(domain)

	return diags
}

func resourceDkimUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceDomainRead(ctx, d, m)
}

func resourceDkimDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*APIClient)

	deleteDkimRequest := api.NewDeleteDKIMKeyRequest()
	items := make([]string, 1)
	items[0] = d.Id()
	deleteDkimRequest.SetItems(items)

	request := c.client.DKIMApi.DeleteDKIMKey(ctx).DeleteDKIMKeyRequest(*deleteDkimRequest)
	_, _, err := c.client.DKIMApi.DeleteDKIMKeyExecute(request)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
