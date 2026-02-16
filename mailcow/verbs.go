package mailcow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/l-with/terraform-provider-mailcow/api"
)

func checkResponse(response api.MailcowResponseArray, resourceName, info string) error {
	t := *response.GetFinalType()
	log.Print("[TRACE] checkResponse t: ", t)
	if t != "success" {
		return errors.New(fmt.Sprintf(
			"%s '%s': %s (%s)",
			resourceName,
			info,
			t,
			*response.GetFinalMsgs(),
		))
	}
	return nil
}

func createRequestSet(mailcowCreateRequest *api.MailcowCreateRequest, res *schema.Resource, data *schema.ResourceData, exclude *[]string, mapArguments *map[string]string) {
	for argument := range (*res).Schema {
		log.Print("[TRACE] createRequestSet argument: ", argument)
		if isElementIn(argument, exclude) {
			log.Print("[TRACE] createRequestSet excluded argument: ", argument)
			continue
		}
		value := data.Get(argument)
		log.Print("[TRACE] createRequestSet set argument: ", getMappedArgument(argument, mapArguments), " := ", value)
		mailcowCreateRequest.Set(getMappedArgument(argument, mapArguments), value)
	}
}

func readRequest(request api.ApiMailcowGetRequest) (map[string]interface{}, error) {
	response, err := request.MailcowExecute()
	if err != nil {
		return nil, err
	}
	result := make(map[string]interface{}, 0)
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func readAllRequest(request api.ApiMailcowGetAllRequest) ([]map[string]interface{}, error) {
	response, err := request.ApiService.MailcowGetAllExecute(request)
	if err != nil {
		return nil, err
	}
	result := make([]map[string]interface{}, 0)
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func updateRequestSetAttr(mailcowUpdateRequest *api.MailcowUpdateRequest, res *schema.Resource, data *schema.ResourceData, exclude *[]string, mapArguments *map[string]string) {
	for argument := range (*res).Schema {
		log.Print("[TRACE] updateRequestSetAttr argument: ", argument)
		if isElementIn(argument, exclude) {
			log.Print("[TRACE] updateRequestSetAttr excluded argument: ", argument)
			continue
		}
		if data.HasChange(argument) {
			log.Print("[TRACE] updateRequestSetAttr changed argument: ", argument)
			mailcowUpdateRequest.SetAttr(getMappedArgument(argument, mapArguments), data.Get(argument))
		}
	}
}

func mailcowCreate(
	ctx context.Context,
	res *schema.Resource,
	d *schema.ResourceData,
	id string,
	exclude *[]string,
	mapArguments *map[string]string,
	mailcowCreateRequest *api.MailcowCreateRequest,
	c *APIClient) error {

	createRequestSet(mailcowCreateRequest, res, d, exclude, mapArguments)

	request := c.client.Api.MailcowCreate(ctx).MailcowCreateRequest(*mailcowCreateRequest)
	response, _, err := c.client.Api.MailcowCreateExecute(request)
	if err != nil {
		return err
	}
	err = checkResponse(response, mailcowCreateRequest.ResourceName, id)
	if err != nil {
		return err
	}
	return nil
}

func mailcowUpdate(
	ctx context.Context,
	res *schema.Resource,
	d *schema.ResourceData,
	exclude *[]string,
	mapArguments *map[string]string,
	mailcowUpdateRequest *api.MailcowUpdateRequest,
	c *APIClient) error {

	updateRequestSetAttr(mailcowUpdateRequest, res, d, exclude, mapArguments)

	mailcowUpdateRequest.SetItem(d.Id())

	response, err := api.MailcowUpdateExecute(ctx, c.client, mailcowUpdateRequest)
	if err != nil {
		return err
	}
	err = checkResponse(response, mailcowUpdateRequest.ResourceName, d.Id())
	if err != nil {
		return err
	}
	return nil
}

func mailcowDelete(
	ctx context.Context,
	d *schema.ResourceData,
	mailcowDeleteRequest *api.MailcowDeleteRequest,
	c *APIClient) (diag.Diagnostics, bool) {

	mailcowDeleteRequest.SetItem(d.Id())

	response, err := api.MailcowDeleteExecute(ctx, c.client, mailcowDeleteRequest)
	if err != nil {
		return diag.FromErr(err), true
	}
	err = checkResponse(response, mailcowDeleteRequest.ResourceName, d.Id())
	if err != nil {
		return diag.FromErr(err), true
	}

	d.SetId("")
	return nil, false
}
