package jsonapitypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	"github.com/rootlyhq/jsonapi"
)

func NewNullableFromString(v types.String) jsonapi.NullableAttr[string] {
	if v.IsUnknown() {
		return jsonapi.NullableAttr[string]{}
	}
	if v.IsNull() {
		return jsonapi.NewNullNullableAttr[string]()
	}
	return jsonapi.NewNullableAttrWithValue(v.ValueString())
}

func NullableStringValue(v jsonapi.NullableAttr[string]) types.String {
	if v, err := v.Get(); err == nil {
		return types.StringValue(v)
	}
	return types.StringNull()
}

func NewNullableFromBool(v types.Bool) jsonapi.NullableAttr[bool] {
	if v.IsUnknown() {
		return jsonapi.NullableAttr[bool]{}
	}
	if v.IsNull() {
		return jsonapi.NewNullNullableAttr[bool]()
	}
	return jsonapi.NewNullableAttrWithValue(v.ValueBool())
}

func NullableBoolValue(v jsonapi.NullableAttr[bool]) types.Bool {
	if v, err := v.Get(); err == nil {
		return types.BoolValue(v)
	}
	return types.BoolNull()
}

func NewNullableFromInt64(v types.Int64) jsonapi.NullableAttr[int64] {
	if v.IsUnknown() {
		return jsonapi.NullableAttr[int64]{}
	}
	if v.IsNull() {
		return jsonapi.NewNullNullableAttr[int64]()
	}
	return jsonapi.NewNullableAttrWithValue(v.ValueInt64())
}

func NullableInt64Value(v jsonapi.NullableAttr[int64]) types.Int64 {
	if v, err := v.Get(); err == nil {
		return types.Int64Value(v)
	}
	return types.Int64Null()
}

func NewNullableFromListOf[T any](ctx context.Context, v supertypes.ListValueOf[T]) (jsonapi.NullableAttr[[]T], diag.Diagnostics) {
	if v.IsUnknown() {
		return jsonapi.NullableAttr[[]T]{}, nil
	}
	if v.IsNull() {
		return jsonapi.NewNullNullableAttr[[]T](), nil
	}

	vv, diags := v.Get(ctx)
	if diags.HasError() {
		return nil, diags
	}
	return jsonapi.NewNullableAttrWithValue(vv), nil
}

func NullableListValueOfSlice[T any](ctx context.Context, v jsonapi.NullableAttr[[]T]) supertypes.ListValueOf[T] {
	if v, err := v.Get(); err == nil {
		return supertypes.NewListValueOfSlice(ctx, v)
	}
	return supertypes.NewListValueOfNull[T](ctx)
}

func NewNullableFromSetOf[T any](ctx context.Context, v supertypes.SetValueOf[T]) (jsonapi.NullableAttr[[]T], diag.Diagnostics) {
	if v.IsUnknown() {
		return jsonapi.NullableAttr[[]T]{}, nil
	}
	if v.IsNull() {
		return jsonapi.NewNullNullableAttr[[]T](), nil
	}

	vv, diags := v.Get(ctx)
	if diags.HasError() {
		return nil, diags
	}
	return jsonapi.NewNullableAttrWithValue(vv), nil
}

func NullableSetValueOfSlice[T any](ctx context.Context, v jsonapi.NullableAttr[[]T]) supertypes.SetValueOf[T] {
	if v, err := v.Get(); err == nil {
		return supertypes.NewSetValueOfSlice(ctx, v)
	}
	return supertypes.NewSetValueOfNull[T](ctx)
}
