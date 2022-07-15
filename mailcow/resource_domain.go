package mailcow

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "github.com/l-with/mailcow-go"
	"io"
	"log"
	"reflect"
	"strconv"
)

func resourceDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainCreate,
		ReadContext:   resourceDomainRead,
		UpdateContext: resourceDomainUpdate,
		DeleteContext: resourceDomainDelete,
		Schema: map[string]*schema.Schema{
			"active": {
				Type:        schema.TypeBool,
				Description: "is domain active or not",
				Default:     true,
				Optional:    true,
			},
			"aliases": {
				Type:        schema.TypeInt,
				Description: "limit count of aliases associated with this domain",
				Default:     400,
				Optional:    true,
			},
			"backupmx": {
				Type:        schema.TypeBool,
				Description: "relay domain or not",
				Default:     false,
				Optional:    true,
			},
			"defquota": {
				Type:        schema.TypeInt,
				Description: "predefined mailbox quota in add mailbox form",
				Default:     3072,
				Optional:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of domain",
				Default:     "mailcow domain",
				Optional:    true,
			},
			"domain": {
				Type:        schema.TypeString,
				Description: "Fully qualified domain name",
				Required:    true,
				ForceNew:    true,
			},
			"gal": {
				Type:        schema.TypeBool,
				Description: "is domain global address list active or not, it enables shared contacts accross domain in SOGo webmail",
				Default:     false,
				Optional:    true,
			},
			"mailboxes": {
				Type:        schema.TypeInt,
				Description: "limit count of mailboxes associated with this domain",
				Default:     10,
				Optional:    true,
			},
			"maxquota": {
				Type:        schema.TypeInt,
				Description: "maximum quota per mailbox",
				Default:     10240,
				Optional:    true,
			},
			"quota": {
				Type:        schema.TypeInt,
				Description: "maximum quota for this domain (for all mailboxes in sum)",
				Default:     10240,
				Optional:    true,
			},
			"restart_sogo": {
				Type:        schema.TypeBool,
				Description: "restart SOGo to activate the domain in SOGo",
				Default:     true,
				Optional:    true,
			},
			"relay_all_recipients": {
				Type:        schema.TypeBool,
				Description: "if not, them you have to create \"dummy\" mailbox for each address to relay",
				Default:     false,
				Optional:    true,
			},
			"relay_unknown_only": {
				Type:        schema.TypeBool,
				Description: "Relay non-existing mailboxes only. Existing mailboxes will be delivered locally.",
				Default:     false,
				Optional:    true,
			},
			"rate_limit": {
				Type:        schema.TypeString,
				Description: "rate limit, decimal with unit s,m,h,d",
				Default:     "10s",
				Optional:    true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*APIClient)

	createDomainRequest := api.NewCreateDomainRequest()
	createDomainRequest.SetActive(d.Get("active").(bool))
	createDomainRequest.SetAliases(float32(d.Get("aliases").(int)))
	createDomainRequest.SetBackupmx(d.Get("backupmx").(bool))
	createDomainRequest.SetDefquota(float32(d.Get("defquota").(int)))
	createDomainRequest.SetDescription(d.Get("description").(string))
	createDomainRequest.SetDomain(d.Get("domain").(string))
	createDomainRequest.SetMailboxes(float32(d.Get("mailboxes").(int)))
	createDomainRequest.SetMaxquota(float32(d.Get("maxquota").(int)))
	createDomainRequest.SetQuota(float32(d.Get("quota").(int)))
	restStartSogo := d.Get("restart_sogo").(bool)
	if restStartSogo {
		createDomainRequest.SetRestartSogo(1.0)
	} else {
		createDomainRequest.SetRestartSogo(0.0)

	}
	createDomainRequest.SetRelayAllRecipients(d.Get("relay_all_recipients").(bool))
	createDomainRequest.SetRelayUnknownOnly(d.Get("relay_unknown_only").(bool))
	iRateLimit, ok := d.GetOk("rate_limit")
	if ok {
		rateLimit := iRateLimit.(string)
		rateFrame := rateLimit[len(rateLimit)-1:]
		createDomainRequest.SetRlFrame(rateFrame)
		rateValue, err := strconv.Atoi(rateLimit[0 : len(rateLimit)-1])
		if err != nil {
			return diag.FromErr(err)
		}
		createDomainRequest.SetRlValue(float32(rateValue))
	}
	tagsInterface := d.Get("tags").([]interface{})
	numTags := len(tagsInterface)
	tags := make([]string, numTags)
	for i, tag := range tagsInterface {
		tags[i] = tag.(string)
	}
	createDomainRequest.SetTags(tags)

	request := c.client.DomainsApi.CreateDomain(ctx).CreateDomainRequest(*createDomainRequest)
	_, _, err := c.client.DomainsApi.CreateDomainExecute(request)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("domain").(string))
	return diags
}

func resourceDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("****** resourceDomainRead ******")

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*APIClient)
	domainName := d.Id()

	request := c.client.DomainsApi.GetDomains(ctx, domainName)

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

	domain := make(map[string]interface{}, 0)
	err = json.NewDecoder(response.Body).Decode(&domain)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("aliases", domain["max_num_aliases_for_domain"])
	if err != nil {
		return diag.FromErr(err)
	}
	domainDefNewMailboxQuota := domain["def_new_mailbox_quota"]
	if domainDefNewMailboxQuota != nil {
		err = d.Set("defquota", int(domainDefNewMailboxQuota.(float64))/(1024*1024))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	err = d.Set("description", domain["description"])
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("domain", domain["domain_name"])
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("mailboxes", domain["max_num_mboxes_for_domain"])
	if err != nil {
		return diag.FromErr(err)
	}
	domainMaxQuotaForMBox := domain["max_quota_for_mbox"]
	if domainMaxQuotaForMBox != nil {
		err = d.Set("maxquota", int(domainMaxQuotaForMBox.(float64))/(1024*1024))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	domainMaxQuotaForDomain := domain["max_quota_for_domain"]
	if domainMaxQuotaForDomain != nil {
		err = d.Set("quota", int(domainMaxQuotaForDomain.(float64))/(1024*1024))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	for _, argumentBool := range []string{
		"active",
		"backupmx",
		"gal",
		"restart_sogo",
		"relay_all_recipients",
		"relay_unknown_only",
	} {
		boolValue := false
		domainArgumentBool := domain[argumentBool]
		if domainArgumentBool != nil {
			if domainArgumentBool.(float64) == 1 {
				boolValue = true
			}
			err = d.Set(argumentBool, boolValue)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	domainQuota := domain["quota"]
	if domainQuota != nil {
		quota, err := strconv.Atoi(fmt.Sprint(domainQuota))
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("quota", quota)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	domainRl := domain["rl"]
	if domainRl != nil {
		rl := make(map[string]string)
		value := reflect.ValueOf(domainRl)
		log.Print("MapKeys: ", value.MapKeys())
		for _, key := range value.MapKeys() {
			rl[fmt.Sprint(key)] = fmt.Sprint(value.MapIndex(key))
		}
		rateLimit := rl["value"] + rl["frame"]

		err = d.Set("rate_limit", rateLimit)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	domainTags := domain["tags"]
	log.Print("domainTags: ", domainTags)
	if domainTags != nil {
		numTags := len(domainTags.([]interface{}))
		tags := make([]string, numTags)
		for i, tag := range domainTags.([]interface{}) {
			tags[i] = tag.(string)
		}
		err = d.Set("tags", tags)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	/*
		else {
			tags := make([]string, 0, 0)
			err = d.Set("tags", tags)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	*/

	d.SetId(domainName)

	return diags
}

func resourceDomainUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("****** resourceDomainUpdate ******")

	c := m.(*APIClient)

	domain := d.Id()

	updateDomainRequest := api.NewUpdateDomainRequest()
	updateDomainRequestAttr := api.NewUpdateDomainRequestAttr()

	if d.HasChange("aliases") {
		updateDomainRequestAttr.SetAliases(float32(d.Get("aliases").(int)))
	}
	if d.HasChange("defquota") {
		updateDomainRequestAttr.SetDefquota(float32(d.Get("defquota").(int)))
	}
	if d.HasChange("description") {
		log.Print("decription has change to: ", d.Get("description").(string))
		updateDomainRequestAttr.SetDescription(d.Get("description").(string))
	}
	if d.HasChange("mailboxes") {
		updateDomainRequestAttr.SetMailboxes(float32(d.Get("mailboxes").(int)))
	}
	if d.HasChange("maxquota") {
		updateDomainRequestAttr.SetMaxquota(float32(d.Get("maxquota").(int)))
	}
	if d.HasChange("quota") {
		updateDomainRequestAttr.SetQuota(float32(d.Get("quota").(int)))
	}
	if d.HasChange("active") {
		updateDomainRequestAttr.SetActive(d.Get("active").(bool))
	}
	if d.HasChange("backupmx") {
		updateDomainRequestAttr.SetBackupmx(d.Get("backupmx").(bool))
	}
	if d.HasChange("gal") {
		updateDomainRequestAttr.SetGal(d.Get("gal").(bool))
	}
	if d.HasChange("relay_all_recipients") {
		updateDomainRequestAttr.SetRelayAllRecipients(d.Get("relay_all_recipients").(bool))
	}
	if d.HasChange("relay_unknown_only") {
		updateDomainRequestAttr.SetRelayUnknownOnly(d.Get("relay_unknown_only").(bool))
	}
	if d.HasChange("rate_limit") {
		iRateLimit := d.Get("rate_limit")
		rateLimit := iRateLimit.(string)
		rateFrame := rateLimit[len(rateLimit)-1:]
		updateDomainRequestAttr.SetRlFrame(rateFrame)
		rateValue, err := strconv.Atoi(rateLimit[0 : len(rateLimit)-1])
		if err != nil {
			return diag.FromErr(err)
		}
		updateDomainRequestAttr.SetRlValue(float32(rateValue))
	}
	if d.HasChange("tags") {
		tagsInterface := d.Get("tags").([]interface{})
		numTags := len(tagsInterface)
		tags := make([]string, numTags)
		for i, tag := range tagsInterface {
			tags[i] = tag.(string)
		}
		updateDomainRequestAttr.SetTags(tags)
	}

	items := make([]string, 1)
	items[0] = domain

	updateDomainRequest.SetItems(items)
	/*
		request := c.client.DomainsApi.CreateDomain(ctx).CreateDomainRequest(*createDomainRequest)
		_, _, err := c.client.DomainsApi.CreateDomainExecute(request)
		if err != nil {
			return diag.FromErr(err)
		}
	*/
	updateDomainRequest.SetAttr(*updateDomainRequestAttr)
	request := c.client.DomainsApi.UpdateDomain(ctx).UpdateDomainRequest(*updateDomainRequest)
	_, _, err := c.client.DomainsApi.UpdateDomainExecute(request)
	if err != nil {
		return diag.FromErr(err)
	}

	/*
		app := resourceStageCaptchaSchemaToProvider(d)

		res, hr, err := c.client.StagesApi.StagesCaptchaUpdate(ctx, d.Id()).CaptchaStageRequest(*app).Execute()
		if err != nil {
			return httpToDiag(d, hr, err)
		}
	*/

	return resourceDomainRead(ctx, d, m)
}

func resourceDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	//return diag.Errorf("not implemented")
	log.Printf("****** resourceDomainDelete ******")

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*APIClient)

	deleteDomainRequest := api.NewDeleteDomainRequest()

	items := make([]string, 1)
	items[0] = d.Id()
	deleteDomainRequest.SetItems(items)

	request := c.client.DomainsApi.DeleteDomain(ctx).DeleteDomainRequest(*deleteDomainRequest)
	_, _, err := c.client.DomainsApi.DeleteDomainExecute(request)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
