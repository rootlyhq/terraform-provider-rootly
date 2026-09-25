package sdkutils

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// FilterToSchema recursively strips keys from data maps that do not exist in the schema.
func FilterToSchema(data map[string]interface{}, blockSchema map[string]*schema.Schema) map[string]interface{} {
	if data == nil || blockSchema == nil {
		return nil
	}

	filtered := make(map[string]interface{})

	for key, val := range data {
		sch, ok := blockSchema[key]
		if !ok || val == nil {
			// Skip keys not declared in the schema or nil values
			continue
		}

		// Handle nested blocks or lists of blocks
		if res, ok := sch.Elem.(*schema.Resource); ok {
			switch v := val.(type) {
			case map[string]interface{}:
				filtered[key] = FilterToSchema(v, res.Schema)
			case []interface{}:
				var filteredList []interface{}
				for _, item := range v {
					if itemMap, ok := item.(map[string]interface{}); ok {
						filteredList = append(filteredList, FilterToSchema(itemMap, res.Schema))
					} else {
						filteredList = append(filteredList, item)
					}
				}
				filtered[key] = filteredList
			default:
				filtered[key] = val
			}
		} else {
			filtered[key] = val
		}
	}

	return filtered
}
