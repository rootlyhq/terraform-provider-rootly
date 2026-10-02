package tools

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// https://github.com/hashicorp/terraform-plugin-sdk/issues/477#issuecomment-1238807249
func EqualIgnoringOrder(key, oldValue, newValue string, d *schema.ResourceData) bool {
	oldArray, newArray, ok := listChangeFromDiffKey(key, d)
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

	return listsAreEqual(oldArray, newArray)
}

func listChangeFromDiffKey(key string, d *schema.ResourceData) ([]interface{}, []interface{}, bool) {
	for {
		key = listPathFromDiffKey(key)
		oldData, newData := d.GetChange(key)
		oldArray, oldOK := oldData.([]interface{})
		newArray, newOK := newData.([]interface{})
		if oldOK && newOK {
			return oldArray, newArray, true
		}

		parent := listPathFromDiffKey(key)
		if parent == key {
			return nil, nil, false
		}
	}
}

func listPathFromDiffKey(key string) string {
	dotIndex := strings.LastIndex(key, ".")
	if dotIndex != -1 {
		return key[:dotIndex]
	}
	return key
}

// toString converts any value to a canonical string representation.
func toString(value interface{}) string {
	// Use JSON marshalling to handle maps, slices, and other structured data.
	// This ensures a consistent representation regardless of type.
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		panic(err) // Handle the error as needed.
	}
	return string(jsonBytes)
}

// listsAreEqual compares two lists of any values, ignoring order.
func listsAreEqual(list1, list2 []interface{}) bool {
	// Convert each value in the lists to a string.
	strList1 := make([]string, len(list1))
	strList2 := make([]string, len(list2))

	for i, value := range list1 {
		strList1[i] = toString(value)
	}
	for i, value := range list2 {
		strList2[i] = toString(value)
	}

	// Sort the string lists.
	sort.Strings(strList1)
	sort.Strings(strList2)

	// Compare the sorted lists.
	if len(strList1) != len(strList2) {
		return false
	}
	for i := range strList1 {
		if strList1[i] != strList2[i] {
			return false
		}
	}
	return true
}

// EqualIgnoringOrderAndFields returns a DiffSuppressFunc that compares lists
// ignoring order and only comparing specified fields. This is useful when
// comparing objects with computed fields that should be ignored.
//
// fieldsToCompare: list of field names to include in comparison (all others ignored)
func EqualIgnoringOrderAndFields(fieldsToCompare []string) schema.SchemaDiffSuppressFunc {
	return func(key, oldValue, newValue string, d *schema.ResourceData) bool {
		oldArray, newArray, ok := listChangeFromDiffKey(key, d)
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

		return listsAreEqualByFields(oldArray, newArray, fieldsToCompare)
	}
}

// filterMapByFields returns a new map containing only the specified fields
func filterMapByFields(m map[string]interface{}, fields []string) map[string]interface{} {
	filtered := make(map[string]interface{})
	for _, field := range fields {
		if value, exists := m[field]; exists {
			filtered[field] = value
		}
	}
	return filtered
}

// listsAreEqualByFields compares two lists of maps, ignoring order and only
// comparing specified fields
func listsAreEqualByFields(list1, list2 []interface{}, fields []string) bool {
	// Convert each value in the lists to a string, filtering by fields
	strList1 := make([]string, len(list1))
	strList2 := make([]string, len(list2))

	for i, value := range list1 {
		if m, ok := value.(map[string]interface{}); ok {
			filtered := filterMapByFields(m, fields)
			strList1[i] = toString(filtered)
		} else {
			strList1[i] = toString(value)
		}
	}
	for i, value := range list2 {
		if m, ok := value.(map[string]interface{}); ok {
			filtered := filterMapByFields(m, fields)
			strList2[i] = toString(filtered)
		} else {
			strList2[i] = toString(value)
		}
	}

	// Sort the string lists.
	sort.Strings(strList1)
	sort.Strings(strList2)

	// Compare the sorted lists.
	if len(strList1) != len(strList2) {
		return false
	}
	for i := range strList1 {
		if strList1[i] != strList2[i] {
			return false
		}
	}
	return true
}
