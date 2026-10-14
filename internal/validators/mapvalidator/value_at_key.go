package mapvalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ validator.Map = ValueStringAtKeyValidator{}

type ValueStringAtKeyValidator struct {
	Key        string
	Validators []validator.String
}

// ValueStringAtKey returns a map validator that runs the given string validators
// against the value at the given key, if the key is present in the map.
func ValueStringAtKey(key string, validators ...validator.String) validator.Map {
	return ValueStringAtKeyValidator{
		Key:        key,
		Validators: validators,
	}
}

func (v ValueStringAtKeyValidator) Description(ctx context.Context) string {
	var descriptions []string
	for _, vv := range v.Validators {
		descriptions = append(descriptions, vv.Description(ctx))
	}
	return fmt.Sprintf("value at key %q must satisfy all validations: %s", v.Key, strings.Join(descriptions, " + "))
}

func (v ValueStringAtKeyValidator) MarkdownDescription(ctx context.Context) string {
	var descriptions []string
	for _, vv := range v.Validators {
		descriptions = append(descriptions, vv.MarkdownDescription(ctx))
	}
	return fmt.Sprintf("value at key `%s` must satisfy all validations: %s", v.Key, strings.Join(descriptions, " + "))
}

func (v ValueStringAtKeyValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	elem, ok := req.ConfigValue.Elements()[v.Key]
	if !ok {
		return
	}

	elemValuable, ok := elem.(basetypes.StringValuable)
	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Validator for Element Type",
			fmt.Sprintf("Map element at key %q is not a string.", v.Key),
		)
		return
	}

	elemValue, diags := elemValuable.ToStringValue(ctx)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	elemReq := validator.StringRequest{
		Path:           req.Path.AtMapKey(v.Key),
		PathExpression: req.PathExpression.AtMapKey(v.Key),
		ConfigValue:    elemValue,
		Config:         req.Config,
	}

	for _, vv := range v.Validators {
		elemResp := &validator.StringResponse{}
		vv.ValidateString(ctx, elemReq, elemResp)
		resp.Diagnostics.Append(elemResp.Diagnostics...)
	}
}
