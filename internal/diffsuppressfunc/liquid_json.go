package diffsuppressfunc

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var _ schema.SchemaDiffSuppressFunc = LiquidJSON

// LiquidJSON suppresses diffs between JSON strings that may contain Liquid.
// Valid JSON is compared semantically. Anything else (e.g. a Liquid template
// that only becomes JSON once rendered) is compared ignoring surrounding
// whitespace, because the API trims it before storing the value.
func LiquidJSON(k, oldValue, newValue string, d *schema.ResourceData) bool {
	var oldJSON, newJSON interface{}
	if json.Unmarshal([]byte(oldValue), &oldJSON) == nil && json.Unmarshal([]byte(newValue), &newJSON) == nil {
		return reflect.DeepEqual(oldJSON, newJSON)
	}
	return strings.TrimSpace(oldValue) == strings.TrimSpace(newValue)
}
