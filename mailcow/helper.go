package mailcow

import "errors"

func ckeckResponse(responseType string, responseMsg []interface{}) error {
	if responseType == "danger" {
		return errors.New("mailcow API response: danger, " + responseMsg[0].(string) + " " + responseMsg[1].(string))
	}
	return nil
}
