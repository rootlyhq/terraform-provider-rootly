package sdkutils

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SDKv2 otherwise treats an empty optional/computed string as unset, retaining state.
func ExplicitEmptyStrings(fields ...string) schema.CustomizeDiffFunc {
	return func(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
		config := diff.GetRawConfig()
		if config.IsNull() || !config.IsKnown() {
			return nil
		}
		for _, field := range fields {
			value := config.GetAttr(field)
			if value.IsNull() || !value.IsKnown() || value.AsString() != "" {
				continue
			}
			if err := diff.SetNew(field, ""); err != nil {
				return err
			}
		}
		return nil
	}
}

// SDKv2 can retain old state for a customized empty-string plan; compare configuration too.
func OptionalComputedStringUpdate(data *schema.ResourceData, field string) *string {
	config := data.GetRawConfig()
	if config.IsNull() || !config.IsKnown() {
		return nil
	}
	value := config.GetAttr(field)
	if value.IsNull() || !value.IsKnown() {
		return nil
	}
	text := value.AsString()
	if !data.HasChange(field) && data.Get(field).(string) == text {
		return nil
	}
	return &text
}
