package api

import (
	"errors"
	"fmt"
)

func (o *MailcowResponseArray) GetId() (error, *string) {
	if !o.HasFinalMsgItem(0) || !o.HasFinalMsgItem(2) {
		return errors.New(fmt.Sprint("msg error: ", o.GetFinalMsgs())), nil
	}
	receipt := *o.GetFinalMsgItem(0)
	if receipt != "alias_added" {
		return errors.New(fmt.Sprint("msg error: ", receipt)), nil
	}
	return nil, o.GetFinalMsgItem(2)
}
