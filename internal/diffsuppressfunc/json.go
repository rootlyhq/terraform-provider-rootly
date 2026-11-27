package diffsuppressfunc

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var _ schema.SchemaDiffSuppressFunc = JSON

// JSON suppresses diffs between semantically equal JSON strings. Values that
// are not valid JSON (e.g. Liquid templates) are compared ignoring surrounding
// whitespace, because the API trims it before storing them.
func JSON(k, oldValue, newValue string, d *schema.ResourceData) bool {
	var oldJSON, newJSON interface{}
	if json.Unmarshal([]byte(oldValue), &oldJSON) == nil && json.Unmarshal([]byte(newValue), &newJSON) == nil {
		return reflect.DeepEqual(oldJSON, newJSON)
	}
	return strings.TrimSpace(oldValue) == strings.TrimSpace(newValue)
}
