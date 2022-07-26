package api

import (
	"context"
	"encoding/json"
	"log"
	"reflect"
)

type MailcowUpdateRequest struct {
	attr         map[string]interface{} `json:"attr,omitempty"`
	items        []string               `json:"items,omitempty"`
	endpoint     string
	ResourceName string
}

func NewUpdateAliasRequest() *MailcowUpdateRequest {
	this := MailcowUpdateRequest{}
	this.attr = make(map[string]interface{})
	this.items = make([]string, 1)
	this.endpoint = "/api/v1/edit/alias"
	this.ResourceName = "resourceAlias"
	return &this
}

func NewUpdateMailboxRequest() *MailcowUpdateRequest {
	this := MailcowUpdateRequest{}
	this.attr = make(map[string]interface{})
	this.items = make([]string, 1)
	this.endpoint = "/api/v1/edit/mailbox"
	this.ResourceName = "resourceMailbox"
	return &this
}

func NewUpdateDomainRequest() *MailcowUpdateRequest {
	this := MailcowUpdateRequest{}
	this.attr = make(map[string]interface{})
	this.items = make([]string, 1)
	this.endpoint = "/api/v1/edit/domain"
	this.ResourceName = "resourceDomain"
	return &this
}

func NewUpdateSyncjobRequest() *MailcowUpdateRequest {
	this := MailcowUpdateRequest{}
	this.attr = make(map[string]interface{})
	this.items = make([]string, 1)
	this.endpoint = "/api/v1/edit/syncjob"
	this.ResourceName = "resourceSyncjob"
	return &this
}

func (o *MailcowUpdateRequest) GetAttr(key string) interface{} {
	if !o.HasAttr(key) {
		var ret bool
		return ret
	}
	return o.attr[key]
}

func (o *MailcowUpdateRequest) GetAttrOk(key string) (interface{}, bool) {
	if !o.HasAttr(key) {
		return nil, false
	}
	return o.GetAttr(key), true
}

func (o *MailcowUpdateRequest) HasAttr(key string) bool {
	if o != nil && o.attr != nil && o.attr[key] != nil {
		return true
	}
	return false
}

func (o *MailcowUpdateRequest) SetAttr(key string, value interface{}) {
	setValue := value
	switch reflect.TypeOf(value).Kind() {
	case reflect.Bool:
		if value.(bool) {
			setValue = 1
		} else {
			setValue = 0
		}
	}
	o.attr[key] = &setValue
}

func (o *MailcowUpdateRequest) GetItem() *string {
	log.Print("[TRACE] GetItem")
	if !o.HasItem() {
		return nil
	}
	item := (*o).items[0]
	return &item
}

func (o *MailcowUpdateRequest) HasItemOk() (*string, bool) {
	if !o.HasItem() {
		var ret string
		return &ret, false
	}
	return o.GetItem(), true
}

func (o *MailcowUpdateRequest) HasItem() bool {
	log.Print("[TRACE] HasItem")
	if o != nil && o.items != nil && len(o.items) != 1 {
		log.Print("[TRACE] HasItem true")
		return true
	}
	log.Print("[TRACE] GetItem false")
	return false
}

func (o *MailcowUpdateRequest) SetItem(v string) {
	o.items[0] = v
}

func (o MailcowUpdateRequest) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.attr != nil {
		toSerialize["attr"] = o.attr
	}
	if o.items != nil {
		toSerialize["items"] = o.items
	}
	return json.Marshal(toSerialize)
}

func MailcowUpdateExecute(ctx context.Context, c *APIClient, mailcowUpdateRequest *MailcowUpdateRequest) (MailcowResponseArray, error) {
	request := c.Api.MailcowUpdate(ctx).MailcowUpdateRequest(*mailcowUpdateRequest)
	response, _, err := c.Api.MailcowUpdateExecute(request)
	return response, err
}
