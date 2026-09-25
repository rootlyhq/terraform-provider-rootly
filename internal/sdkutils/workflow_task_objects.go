package sdkutils

import (
	"fmt"
	"maps"
)

type WorkflowTaskObjectFields map[string]WorkflowTaskObjectFields

func WorkflowTaskBlocksToObjects(params map[string]interface{}, fields WorkflowTaskObjectFields) (map[string]interface{}, error) {
	if params == nil {
		return nil, fmt.Errorf("workflow task parameters must be an object")
	}
	result := maps.Clone(params)
	for name, nestedFields := range fields {
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
		object, err := WorkflowTaskBlocksToObjects(values, nestedFields)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		result[name] = object
	}
	return result, nil
}

func WorkflowTaskObjectsToBlocks(params map[string]interface{}, fields WorkflowTaskObjectFields) (map[string]interface{}, error) {
	if params == nil {
		return nil, fmt.Errorf("workflow task parameters must be an object")
	}
	result := maps.Clone(params)
	for name, nestedFields := range fields {
		if params[name] == nil {
			result[name] = []interface{}{}
			continue
		}
		values, ok := params[name].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%s must be an object", name)
		}
		block, err := WorkflowTaskObjectsToBlocks(values, nestedFields)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		result[name] = []interface{}{block}
	}
	return result, nil
}
