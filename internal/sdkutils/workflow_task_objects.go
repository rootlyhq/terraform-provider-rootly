package sdkutils

import (
	"fmt"
	"maps"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func FlattenWorkflowTaskObjects(params map[string]interface{}, fields map[string]*schema.Schema) (map[string]interface{}, error) {
	if params == nil {
		return nil, fmt.Errorf("workflow task parameters must be an object")
	}
	result := maps.Clone(params)
	for name, field := range fields {
		object, ok := workflowTaskObjectSchema(field)
		if !ok {
			continue
		}
		if params[name] == nil {
			result[name] = nil
			continue
		}
		blocks, ok := params[name].([]interface{})
		if !ok || len(blocks) > 1 {
			return nil, fmt.Errorf("%s must contain one object block", name)
		}
		if len(blocks) == 0 {
			result[name] = nil
			continue
		}
		values, ok := blocks[0].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%s must contain an object", name)
		}
		flattened, err := FlattenWorkflowTaskObjects(values, object.Schema)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		result[name] = flattened
	}
	return result, nil
}

func ExpandWorkflowTaskObjects(params map[string]interface{}, fields map[string]*schema.Schema) (map[string]interface{}, error) {
	if params == nil {
		return nil, fmt.Errorf("workflow task parameters must be an object")
	}
	result := maps.Clone(params)
	for name, field := range fields {
		object, ok := workflowTaskObjectSchema(field)
		if !ok {
			continue
		}
		if params[name] == nil {
			result[name] = []interface{}{}
			continue
		}
		values, ok := params[name].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%s must be an object", name)
		}
		expanded, err := ExpandWorkflowTaskObjects(values, object.Schema)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		result[name] = []interface{}{expanded}
	}
	return result, nil
}

func workflowTaskObjectSchema(field *schema.Schema) (*schema.Resource, bool) {
	object, ok := field.Elem.(*schema.Resource)
	return object, ok && field.Type == schema.TypeList && field.MaxItems == 1
}
