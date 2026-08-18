package jsonapitypes

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/rootlyhq/jsonapi"
)

func NullableStringValue(v jsonapi.NullableAttr[string]) types.String {
	if v, err := v.Get(); err == nil {
		return types.StringValue(v)
	}
	return types.StringNull()
}

func NullableBoolValue(v jsonapi.NullableAttr[bool]) types.Bool {
	if v, err := v.Get(); err == nil {
		return types.BoolValue(v)
	}
	return types.BoolNull()
}

func NullableInt64Value(v jsonapi.NullableAttr[int64]) types.Int64 {
	if v, err := v.Get(); err == nil {
		return types.Int64Value(v)
	}
	return types.Int64Null()
}
