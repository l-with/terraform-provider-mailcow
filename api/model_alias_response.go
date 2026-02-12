package api

import (
	"errors"
	"fmt"
)

func (o *MailcowResponseArray) GetAliasId() (*string, error) {
	if !o.HasFinalMsgItem(0) || !o.HasFinalMsgItem(2) {
		return nil, errors.New(fmt.Sprint("msg error: ", o.GetFinalMsgs()))
	}
	receipt := *o.GetFinalMsgItem(0)
	if receipt != "alias_added" {
		return nil, errors.New(fmt.Sprint("msg error: ", receipt))
	}
	return o.GetFinalMsgItem(2), nil
}

func (o *MailcowResponseArray) GetAliasDomainId() (*string, error) {
	if !o.HasFinalMsgItem(0) || !o.HasFinalMsgItem(1) {
		return nil, errors.New(fmt.Sprint("msg error: ", o.GetFinalMsgs()))
	}
	receipt := *o.GetFinalMsgItem(0)
	if receipt != "aliasd_added" {
		return nil, errors.New(fmt.Sprint("msg error: ", receipt))
	}
	return o.GetFinalMsgItem(1), nil
}
