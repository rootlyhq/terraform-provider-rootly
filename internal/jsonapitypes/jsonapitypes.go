package jsonapitypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
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

func NullableListOfValue[T any](ctx context.Context, v jsonapi.NullableAttr[[]T]) supertypes.ListValueOf[T] {
	if v, err := v.Get(); err == nil {
		return supertypes.NewListValueOfSlice(ctx, v)
	}
	return supertypes.NewListValueOfNull[T](ctx)
}
