package api

import (
	"fmt"
	"log"
	"reflect"
)

type MailcowResponseArray []struct {
	Type *string `json:"type,omitempty"`
	// contains request object
	Log []interface{} `json:"log,omitempty"`
	Msg interface{}   `json:"msg,omitempty"`
}

type MailcowResponse struct {
	Log  []interface{} `json:"log,omitempty"`
	Msg  []interface{} `json:"msg,omitempty"`
	Type *string       `json:"type,omitempty"`
}

func NewMailcowResponse() *MailcowResponseArray {
	this := MailcowResponseArray{}
	return &this
}

func (o *MailcowResponseArray) IsSlice() bool {
	if o == nil {
		return false
	}
	log.Print("[TRACE] IsSlice Kind: ", reflect.TypeOf(*o).Kind())
	return reflect.TypeOf(*o).Kind() == reflect.Slice
}

func (o *MailcowResponseArray) GetFinalType() *string {
	if !o.HasFinalType() {
		var ret string
		return &ret
	}
	return (*o)[len(*o)-1].Type
}

func (o *MailcowResponseArray) GetFinalTypeOk() (*string, bool) {
	if !o.HasFinalType() {
		return nil, false
	}
	return o.GetFinalType(), true
}

func (o *MailcowResponseArray) HasFinalType() bool {
	if o.IsSlice() && len(*o) > 0 && (*o)[len(*o)-1].Type != nil {
		return true
	}
	return false
}

func (o *MailcowResponseArray) GetFinalLog() []interface{} {
	if !o.HasFinalLog() {
		var ret []interface{}
		return ret
	}
	return (*o)[len(*o)-1].Log
}

func (o *MailcowResponseArray) GetFinalLogOk() ([]interface{}, bool) {
	if !o.HasFinalLog() {
		return nil, false
	}
	return o.GetFinalLog(), true
}

func (o *MailcowResponseArray) HasFinalLog() bool {
	if o.IsSlice() && (*o)[len(*o)-1].Log != nil {
		return true
	}
	return false
}

func (o *MailcowResponseArray) HasFinalMsg() bool {
	log.Print("[TRACE] HasFinalMsg len: ", len(*o))
	if o.IsSlice() && len(*o) > 0 && (*o)[len(*o)-1].Msg != nil {
		return true
	}
	log.Print("[TRACE] HasFinalMsg return false ")
	return false
}

func (o *MailcowResponseArray) GetFinalMsg() interface{} {
	if !o.HasFinalMsg() {
		return nil
	}
	return (*o)[len(*o)-1].Msg
}

func (o *MailcowResponseArray) GetFinalMsgs() *string {
	if !o.HasFinalMsgs() {
		var ret string
		return &ret
	}
	mailcowResponseMsg := o.GetFinalMsg()
	msgs := ""
	if reflect.ValueOf(mailcowResponseMsg).Kind() == reflect.String {
		msgs = mailcowResponseMsg.(string)
		return &msgs
	}
	for _, msgItem := range mailcowResponseMsg.([]interface{}) {
		if msgs != "" {
			msgs = msgs + ", "
		}
		msgs = msgs + fmt.Sprint(msgItem)
	}
	return &msgs
}

func (o *MailcowResponseArray) GetFinalMsgsOk() (*string, bool) {
	if !o.HasFinalMsgs() {
		return nil, false
	}
	return o.GetFinalMsgs(), true
}

func (o *MailcowResponseArray) HasFinalMsgs() bool {
	if o.IsSlice() && o.HasFinalMsg() {
		return true
	}
	return false
}

func (o *MailcowResponseArray) GetFinalMsgItem(i int) *string {
	if !o.HasFinalMsgItem(i) {
		var ret string
		return &ret
	}
	mailcowMsgs := o.GetFinalMsg().([]interface{})
	msgItem := fmt.Sprint(mailcowMsgs[i])
	return &msgItem
}

func (o *MailcowResponseArray) HasFinalMsgItem(i int) bool {
	if i < 0 || o == nil || !o.IsSlice() || o.GetFinalMsg() == nil {
		return false
	}
	mailcowResponseMsg := o.GetFinalMsg()
	if reflect.ValueOf(mailcowResponseMsg).Kind() == reflect.String {
		return false
	}
	mailcowMsgs := mailcowResponseMsg.([]interface{})
	if i >= len(mailcowMsgs) {
		return false
	}
	return true
}
