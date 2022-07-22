package mailcow

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
)

func checkResponse(response []map[string]interface{}, funcName interface{}, info string) error {
	if response[len(response)-1]["type"].(string) != "success" {
		return errors.New(fmt.Sprintf(
			"%s '%s': %s - %s",
			runtime.FuncForPC(reflect.ValueOf(funcName).Pointer()).Name(),
			info,
			response[len(response)-1]["type"].(string),
			getMsg(response),
		))
	}
	return nil
}

func getMsg(response []map[string]interface{}) string {
	responseMsg := response[len(response)-1]["msg"]
	if reflect.ValueOf(responseMsg).Kind() == reflect.String {
		return responseMsg.(string)
	}
	msg := ""
	for _, msgItem := range responseMsg.([]interface{}) {
		msg = msg + ", " + msgItem.(string)
	}
	return msg
}
