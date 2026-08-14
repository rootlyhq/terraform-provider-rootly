package setvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ validator.Set = UniqueByAttributeValidator{}

type UniqueByAttributeValidator struct {
	AttributeName string
}

// UniqueByAttribute returns a set validator that checks that a given attribute
// is unique across all elements in the set.
func UniqueByAttribute(attributeName string) validator.Set {
	return UniqueByAttributeValidator{
		AttributeName: attributeName,
	}
}

func (v UniqueByAttributeValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("Ensures that all elements in the set have a unique value for %s.", v.AttributeName)
}

func (v UniqueByAttributeValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Ensures that all elements in the set have a unique value for `%s`.", v.AttributeName)
}

func (v UniqueByAttributeValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	elements := req.ConfigValue.Elements()
	// Track the first path where each attribute value was encountered
	seenValues := make(map[attr.Value]path.Path)

	for _, elem := range elements {
		var attrs map[string]attr.Value

		if objVal, ok := elem.(basetypes.ObjectValuable); ok {
			obj, diags := objVal.ToObjectValue(ctx)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			attrs = obj.Attributes()
		} else {
			continue
		}

		// Exact path to the target attribute on this specific set element
		elemAttrPath := req.Path.AtSetValue(elem).AtName(v.AttributeName)

		attrValue, exists := attrs[v.AttributeName]
		if !exists {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Attribute Validator Target",
				fmt.Sprintf("Attribute %q does not exist on element at %s.", v.AttributeName, req.Path.String()),
			)
			return
		}

		if attrValue.IsUnknown() || attrValue.IsNull() {
			continue
		}

		if firstSeenPath, found := seenValues[attrValue]; found {
			resp.Diagnostics.AddAttributeError(
				elemAttrPath,
				"Duplicate Attribute Value",
				fmt.Sprintf(
					"Attribute %q has duplicate value %s.\n\n"+
						"• Conflicting Path 1: %s\n"+
						"• Conflicting Path 2: %s\n\n"+
						"Values for %q must be unique across all elements in the set.",
					v.AttributeName,
					attrValue.String(),
					firstSeenPath.String(),
					elemAttrPath.String(),
					v.AttributeName,
				),
			)
			return
		}

		seenValues[attrValue] = elemAttrPath
	}
}
