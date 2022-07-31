package api

import (
	"context"
	"encoding/json"
	"log"
)

type MailcowDeleteRequest struct {
	item         *string
	endpoint     string
	ResourceName string
}

func NewDeleteAliasRequest() *MailcowDeleteRequest {
	this := MailcowDeleteRequest{}
	this.endpoint = "/api/v1/delete/alias"
	this.ResourceName = "resourceAlias"
	return &this
}

func NewDeleteDomainRequest() *MailcowDeleteRequest {
	this := MailcowDeleteRequest{}
	this.endpoint = "/api/v1/delete/domain"
	this.ResourceName = "resourceDomain"
	return &this
}

func NewDeleteMailboxRequest() *MailcowDeleteRequest {
	this := MailcowDeleteRequest{}
	this.endpoint = "/api/v1/delete/mailbox"
	this.ResourceName = "resourceMailbox"
	return &this
}

func NewDeleteDkimRequest() *MailcowDeleteRequest {
	this := MailcowDeleteRequest{}
	this.endpoint = "/api/v1/delete/dkim"
	this.ResourceName = "resourceDkim"
	return &this
}

func NewDeleteSyncjobRequest() *MailcowDeleteRequest {
	this := MailcowDeleteRequest{}
	this.endpoint = "/api/v1/delete/syncjob"
	this.ResourceName = "resourceSyncjob"
	return &this
}

func NewDeleteOAuth2ClientRequest() *MailcowDeleteRequest {
	this := MailcowDeleteRequest{}
	this.endpoint = "/api/v1/delete/oauth2-client"
	this.ResourceName = "resourceOAuth2Client"
	return &this
}

func (o *MailcowDeleteRequest) GetItem() *string {
	log.Print("[TRACE] GetItem")
	if !o.HasItem() {
		return nil
	}
	return o.item
}

func (o *MailcowDeleteRequest) HasItemOk() (*string, bool) {
	if !o.HasItem() {
		var ret string
		return &ret, false
	}
	return o.GetItem(), true
}

func (o *MailcowDeleteRequest) HasItem() bool {
	log.Print("[TRACE] HasItem")
	if o != nil && o.item != nil {
		log.Print("[TRACE] HasItem true")
		return true
	}
	log.Print("[TRACE] HasItem false")
	return false
}

func (o *MailcowDeleteRequest) SetItem(v string) {
	o.item = &v
}

func (o MailcowDeleteRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.GetItem())
}

func MailcowDeleteExecute(ctx context.Context, c *APIClient, mailcowDeleteRequest *MailcowDeleteRequest) (MailcowResponseArray, error) {
	request := c.Api.MailcowDelete(ctx).MailcowDeleteRequest(*mailcowDeleteRequest)
	response, _, err := c.Api.MailcowDeleteExecute(request)
	return response, err
}
