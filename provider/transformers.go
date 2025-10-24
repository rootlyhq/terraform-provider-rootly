package provider

import (
	"reflect"

	"github.com/buxizhizhoum/inflection"
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
