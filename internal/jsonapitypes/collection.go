package jsonapitypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	"github.com/rootlyhq/jsonapi"
)

// nestedContainer is an interface that represents a nested container.
// This is satisfied by supertypes.ListNestedObjectValueOf[U] and supertypes.SetNestedObjectValueOf[U].
type nestedContainer[U any] interface {
	IsUnknown() bool
	IsNull() bool
	Get(ctx context.Context) ([]*U, diag.Diagnostics)
}

type NullableListOutcome int

const (
	OutcomeOmit NullableListOutcome = iota
	OutcomeNull
	OutcomeEmptyList
)

type NullableListConfig struct {
	OnUnknown NullableListOutcome
	OnNull    NullableListOutcome
	OnEmpty   NullableListOutcome
}

var DefaultNullableListConfig = NullableListConfig{
	OnUnknown: OutcomeOmit,
	OnNull:    OutcomeEmptyList,
	OnEmpty:   OutcomeEmptyList,
}

func resolveOutcome[T any](o NullableListOutcome) jsonapi.NullableAttr[[]T] {
	switch o {
	case OutcomeNull:
		return jsonapi.NewNullNullableAttr[[]T]()
	case OutcomeEmptyList:
		return jsonapi.NewNullableAttrWithValue([]T{})
	default: // OutcomeOmit
		return jsonapi.NullableAttr[[]T]{}
	}
}

func convertNullable[T any, U any, C nestedContainer[U]](
	ctx context.Context,
	container C,
	cfg NullableListConfig,
	convert func(ctx context.Context, item *U) (*T, diag.Diagnostics),
) (jsonapi.NullableAttr[[]T], diag.Diagnostics) {
	var diags diag.Diagnostics

	if container.IsUnknown() {
		return resolveOutcome[T](cfg.OnUnknown), diags
	}
	if container.IsNull() {
		return resolveOutcome[T](cfg.OnNull), diags
	}

	items, d := container.Get(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return jsonapi.NullableAttr[[]T]{}, diags
	}
	if len(items) == 0 {
		return resolveOutcome[T](cfg.OnEmpty), diags
	}

	result := make([]T, 0, len(items))
	for i, item := range items {
		if item == nil {
			diags.AddError("null value", fmt.Sprintf("cannot convert null item at index %d", i))
			continue
		}
		converted, d := convert(ctx, item)
		diags.Append(d...)
		if diags.HasError() {
			continue
		}
		if converted == nil {
			diags.AddError("null value", fmt.Sprintf("cannot convert null item at index %d", i))
			continue
		}
		result = append(result, *converted)
	}

	if diags.HasError() {
		return jsonapi.NullableAttr[[]T]{}, diags
	}
	return jsonapi.NewNullableAttrWithValue(result), diags
}

func ConvertNullableList[T any, U any](
	ctx context.Context,
	list supertypes.ListNestedObjectValueOf[U],
	cfg NullableListConfig,
	convert func(ctx context.Context, item *U) (*T, diag.Diagnostics),
) (jsonapi.NullableAttr[[]T], diag.Diagnostics) {
	return convertNullable(ctx, list, cfg, convert)
}

func ConvertNullableSet[T any, U any](
	ctx context.Context,
	set supertypes.SetNestedObjectValueOf[U],
	cfg NullableListConfig,
	convert func(ctx context.Context, item *U) (*T, diag.Diagnostics),
) (jsonapi.NullableAttr[[]T], diag.Diagnostics) {
	return convertNullable(ctx, set, cfg, convert)
}

func convertToNestedModel[T any, U any, C any](
	ctx context.Context,
	apiAttr jsonapi.NullableAttr[[]U],
	diags *diag.Diagnostics,
	fromApi func(ctx context.Context, item *T, apiItem U) diag.Diagnostics,
	newNull func(ctx context.Context) C,
	newFromSlice func(ctx context.Context, items []T) C,
) C {
	v, err := apiAttr.Get()
	if err != nil {
		return newNull(ctx)
	}

	items := make([]T, len(v))
	for i, vv := range v {
		diags.Append(fromApi(ctx, &items[i], vv)...)
	}
	return newFromSlice(ctx, items)
}

func ConvertToListModel[T any, U any](
	ctx context.Context,
	apiAttr jsonapi.NullableAttr[[]U],
	diags *diag.Diagnostics,
	fromApi func(ctx context.Context, item *T, apiItem U) diag.Diagnostics,
) supertypes.ListNestedObjectValueOf[T] {
	return convertToNestedModel[T, U, supertypes.ListNestedObjectValueOf[T]](
		ctx, apiAttr, diags, fromApi,
		supertypes.NewListNestedObjectValueOfNull[T],
		supertypes.NewListNestedObjectValueOfValueSlice[T],
	)
}

func ConvertToSetModel[T any, U any](
	ctx context.Context,
	apiAttr jsonapi.NullableAttr[[]U],
	diags *diag.Diagnostics,
	fromApi func(ctx context.Context, item *T, apiItem U) diag.Diagnostics,
) supertypes.SetNestedObjectValueOf[T] {
	return convertToNestedModel[T, U, supertypes.SetNestedObjectValueOf[T]](
		ctx, apiAttr, diags, fromApi,
		supertypes.NewSetNestedObjectValueOfNull[T],
		supertypes.NewSetNestedObjectValueOfValueSlice[T],
	)
}
