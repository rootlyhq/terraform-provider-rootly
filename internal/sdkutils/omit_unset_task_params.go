package sdkutils

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// OmitUnsetTaskParams drops task params the configuration leaves unset, so the
// API applies its own default instead of receiving the zero value.
func OmitUnsetTaskParams(data *schema.ResourceData, params map[string]interface{}, keys ...string) {
	config := data.GetRawConfig()
	if config.IsNull() || !config.IsKnown() {
		return
	}
	taskParams := config.GetAttr("task_params")
	var block cty.Value
	if !taskParams.IsNull() && taskParams.IsKnown() && taskParams.LengthInt() > 0 {
		block = taskParams.Index(cty.NumberIntVal(0))
	}
	for _, key := range keys {
		if block == cty.NilVal || block.IsNull() || !block.Type().HasAttribute(key) || block.GetAttr(key).IsNull() {
			delete(params, key)
		}
	}
}
