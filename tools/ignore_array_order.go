package tools

import (
	"fmt"
	"sort"
	"strings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"reflect"
)

// https://github.com/hashicorp/terraform-plugin-sdk/issues/477#issuecomment-1238807249
func EqualIgnoringOrder(key, oldValue, newValue string, d *schema.ResourceData) bool {
	// The key is a path not the list itself, e.g. "events.0"
 	lastDotIndex := strings.LastIndex(key, ".")
 	if lastDotIndex != -1 {
 		key = string(key[:lastDotIndex])
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

 	// Convert data to string arrays
 	oldItems := make([]string, len(oldArray))
 	newItems := make([]string, len(newArray))
 	for i, oldItem := range oldArray {
 		oldItems[i] = fmt.Sprint(oldItem)
 	}
 	for j, newItem := range newArray {
 		newItems[j] = fmt.Sprint(newItem)
 	}
 	// Ensure consistent sorting before comparison, to suppress unnecessary change detections
 	sort.Strings(oldItems)
 	sort.Strings(newItems)
 	return reflect.DeepEqual(oldItems, newItems)
}
