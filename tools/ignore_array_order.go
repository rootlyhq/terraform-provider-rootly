package tools

import (
	"encoding/json"
	"sort"
	"strings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// https://github.com/hashicorp/terraform-plugin-sdk/issues/477#issuecomment-1238807249
func EqualIgnoringOrder(key, oldValue, newValue string, d *schema.ResourceData) bool {
	// For a list, the key is path to the element, rather than the list.
	// E.g. "node_groups.2.ips.0"
	dotIndex := strings.Index(key, ".")
	if dotIndex != -1 {
		key = string(key[:dotIndex])
	}

	oldData, newData := d.GetChange(key)
	if oldData == nil || newData == nil {
		return false
	}

 	oldArray := oldData.([]interface{})
 	newArray := newData.([]interface{})
 	if len(oldArray) != len(newArray) {
 		// Items added or removed, always detect as changed
 		return false
 	}

	return listsAreEqual(oldArray, newArray)
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
