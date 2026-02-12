package mailcow

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDataSet(rd *schema.ResourceData, argument string, value any, elem *schema.Schema) error {
	stringValue := fmt.Sprint(value)
	log.Printf("[TRACE] resourceDataSet %s expected type %s, value %s", argument, elem.Type, stringValue)
	setValue := value
	var err error
	switch elem.Type {
	case schema.TypeBool:
		if stringValue == "1" {
			setValue = true
		} else if stringValue == "0" {
			setValue = false
		} else if stringValue == "" {
			setValue = false
		}
	case schema.TypeInt:
		var setValueInt int
		setValueInt, err = strconv.Atoi(stringValue)
		if err != nil {
			return err
		}
		setValue = setValueInt
	case schema.TypeList:
		if value != nil {
			num := len(value.([]interface{}))
			list := make([]string, num)
			for i, item := range value.([]interface{}) {
				list[i] = item.(string)
			}
			setValue = list
		} else {
			list := make([]string, 0)
			setValue = list
		}
	case schema.TypeString:
		setValue = stringValue
	default:
		setValue = stringValue
	}
	log.Printf("[TRACE] resourceDataSet %s setVvalue %s", argument, setValue)
	return rd.Set(argument, setValue)
}

func setResourceData(res *schema.Resource, data *schema.ResourceData, resource *map[string]interface{}, exclude *[]string, only *[]string) error {
	var err error
	for argument, elem := range (*res).Schema {
		if isElementIn(argument, exclude) {
			continue
		}
		if only != nil && !isElementIn(argument, only) {
			continue
		}
		err = resourceDataSet(data, argument, (*resource)[argument], elem)
		if err != nil {
			return err
		}
	}
	return nil
}

func getMappedArgument(argument string, mapArguments *map[string]string) string {
	if mapArguments == nil {
		return argument
	}
	mapppedArgument := argument
	if value, ok := (*mapArguments)[argument]; ok {
		mapppedArgument = value
	}
	return mapppedArgument
}

func isElementIn(argument string, arguments *[]string) bool {
	if arguments == nil {
		return false
	}
	for _, elem := range *arguments {
		if argument == elem {
			return true
		}
	}
	return false
}

func randomLowerCaseString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}
