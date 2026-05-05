package provider

import (
	"reflect"

	"github.com/buxizhizhoum/inflection"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func structToLowerFirstMap(in interface{}) map[string]interface{} {
	v := reflect.ValueOf(in)
	vType := v.Type()

	result := make(map[string]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		name := vType.Field(i).Name
		result[inflection.Underscore(name)] = v.Field(i).Interface()
	}

	return result
}

func filterMapKeys(m map[string]interface{}, allowed map[string]*schema.Schema) map[string]interface{} {
	result := make(map[string]interface{}, len(allowed))
	for k, v := range m {
		if _, ok := allowed[k]; ok {
			result[k] = v
		}
	}
	return result
}
