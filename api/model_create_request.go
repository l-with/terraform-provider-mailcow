package api

import (
	"encoding/json"
	"log"
	"reflect"
)

type MailcowCreateRequest struct {
	payload      map[string]interface{} `json:"payload,omitempty"`
	attr         map[string]interface{} `json:"attr,omitempty"`
	endpoint     string
	ResourceName string
}

func NewCreateAliasRequest() *MailcowCreateRequest {
	this := MailcowCreateRequest{}
	this.payload = make(map[string]interface{})
	this.endpoint = "/api/v1/add/alias"
	this.ResourceName = "resourceAlias"
	return &this
}

func NewCreateDomainRequest() *MailcowCreateRequest {
	this := MailcowCreateRequest{}
	this.payload = make(map[string]interface{})
	this.endpoint = "/api/v1/add/domain"
	this.ResourceName = "resourceAlias"
	return &this
}

func NewCreateMailboxRequest() *MailcowCreateRequest {
	this := MailcowCreateRequest{}
	this.payload = make(map[string]interface{})
	this.endpoint = "/api/v1/add/mailbox"
	this.ResourceName = "resourceAlias"
	return &this
}

func NewCreateDkimRequest() *MailcowCreateRequest {
	this := MailcowCreateRequest{}
	this.payload = make(map[string]interface{})
	this.endpoint = "/api/v1/add/dkim"
	this.ResourceName = "resourceAlias"
	return &this
}

func NewCreateSyncjobRequest() *MailcowCreateRequest {
	this := MailcowCreateRequest{}
	this.payload = make(map[string]interface{})
	this.endpoint = "/api/v1/add/syncjob"
	this.ResourceName = "resourceAlias"
	return &this
}

func NewCreateOAuth2ClientRequest() *MailcowCreateRequest {
	this := MailcowCreateRequest{}
	this.payload = make(map[string]interface{})
	this.endpoint = "/api/v1/add/oauth2-client"
	this.ResourceName = "resourceOAuth2Client"
	return &this
}

func (o *MailcowCreateRequest) Get(key string) interface{} {
	if !o.Has(key) {
		var ret bool
		return ret
	}
	return o.payload[key]
}

func (o *MailcowCreateRequest) GetOk(key string) (interface{}, bool) {
	if !o.Has(key) {
		return nil, false
	}
	return o.Get(key), true
}

func (o *MailcowCreateRequest) Has(key string) bool {
	if o != nil && o.payload != nil && o.payload[key] != nil {
		return true
	}
	return false
}

func (o *MailcowCreateRequest) Set(key string, value interface{}) {
	setValue := value
	switch reflect.TypeOf(value).Kind() {
	case reflect.Bool:
		if value.(bool) {
			setValue = 1
		} else {
			setValue = 0
		}
	}
	log.Print("[TRACE] CreateRequest Set key: ", key, ", value: ", setValue)
	o.payload[key] = &setValue
}

func (o *MailcowCreateRequest) GetAttr(key string) interface{} {
	if !o.HasAttr(key) {
		var ret bool
		return ret
	}
	return o.attr[key]
}

func (o *MailcowCreateRequest) GetAttrOk(key string) (interface{}, bool) {
	if !o.HasAttr(key) {
		return nil, false
	}
	return o.GetAttr(key), true
}

func (o *MailcowCreateRequest) HasAttr(key string) bool {
	if o != nil && o.attr != nil && o.attr[key] != nil {
		return true
	}
	return false
}

func (o *MailcowCreateRequest) SetAttr(key string, value interface{}) {
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

func (o MailcowCreateRequest) MarshalJSON(requestSpec map[string]interface{}) ([]byte, error) {
	toSerialize := map[string]interface{}{}
	//if o.attr != nil {
	//	toSerialize["attr"] = o.attr
	//}
	for key := range requestSpec {
		//key := element.(map)
		if o.payload[key] != nil {
			toSerialize[key] = o.payload[key]
		}
	}
	return json.Marshal(toSerialize)
}
