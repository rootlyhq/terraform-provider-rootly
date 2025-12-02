package rootlytypes

import "github.com/hashicorp/terraform-plugin-framework/types/basetypes"

type Int = basetypes.Int64Value

func IntNull() basetypes.Int64Value {
	return basetypes.NewInt64Null()
}

func IntUnknown() basetypes.Int64Value {
	return basetypes.NewInt64Unknown()
}

func IntValue(value int) basetypes.Int64Value {
	return basetypes.NewInt64Value(int64(value))
}

func IntPointerValue(value *int) basetypes.Int64Value {
	if value == nil {
		return IntNull()
	}
	return IntValue(*value)
}
