package diffsuppressfunc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v5/tools"
)

// EqualIgnoringOrderOmittingFields returns a DiffSuppressFunc that compares
// lists of objects ignoring order, after dropping from each object the fields
// returned by omit. Use it for fields the API rewrites on some items, so a
// configured value that the server never keeps does not cause a perpetual diff.
func EqualIgnoringOrderOmittingFields(omit func(item map[string]interface{}) []string) schema.SchemaDiffSuppressFunc {
	return func(key, oldValue, newValue string, d *schema.ResourceData) bool {
		oldArray, newArray, ok := tools.ListChangeFromDiffKey(key, d)
		if !ok {
			return false
		}
		if len(oldArray) != len(newArray) {
			// Items added or removed, always detect as changed
			return false
		}

		// Workaround to detect lists being removed from plan
		if len(oldArray) > 0 && oldValue != newValue && newValue == "0" && oldArray[0] != "0" {
			return false
		}

		return tools.ListsAreEqualIgnoringOrder(omitFieldsFromList(oldArray, omit), omitFieldsFromList(newArray, omit))
	}
}

// omitFieldsFromList returns a copy of list where each map item has the
// fields returned by omit removed. Non-map items are kept as they are.
func omitFieldsFromList(list []interface{}, omit func(item map[string]interface{}) []string) []interface{} {
	result := make([]interface{}, len(list))
	for i, value := range list {
		m, ok := value.(map[string]interface{})
		if !ok {
			result[i] = value
			continue
		}
		fields := omit(m)
		if len(fields) == 0 {
			result[i] = m
			continue
		}
		filtered := make(map[string]interface{}, len(m))
		for k, v := range m {
			filtered[k] = v
		}
		for _, field := range fields {
			delete(filtered, field)
		}
		result[i] = filtered
	}
	return result
}
