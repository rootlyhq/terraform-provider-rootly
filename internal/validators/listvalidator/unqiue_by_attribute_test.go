package listvalidator_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/rootlyhq/terraform-provider-rootly/v5/internal/validators/listvalidator"
)

func memberObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"member_id":   types.StringType,
		"member_type": types.StringType,
		"position":    types.Int64Type,
	}
}

func createMemberObject(memberID, memberType string, position *int64) types.Object {
	posAttr := types.Int64Null()
	if position != nil {
		posAttr = types.Int64Value(*position)
	}

	obj, _ := types.ObjectValue(
		memberObjectType(),
		map[string]attr.Value{
			"member_id":   types.StringValue(memberID),
			"member_type": types.StringValue(memberType),
			"position":    posAttr,
		},
	)
	return obj
}

func TestUniqueByAttributeValidator(t *testing.T) {
	t.Parallel()

	intPtr := func(v int64) *int64 { return &v }

	tests := map[string]struct {
		configVal    basetypes.ListValue
		expectedErr  bool
		expectedDiag string
	}{
		"valid list - unique positions": {
			configVal: types.ListValueMust(
				types.ObjectType{AttrTypes: memberObjectType()},
				[]attr.Value{
					createMemberObject("usr-1", "User", intPtr(1)),
					createMemberObject("usr-2", "User", intPtr(2)),
					createMemberObject("sch-1", "Schedule", intPtr(3)),
				},
			),
			expectedErr: false,
		},
		"valid list - null positions ignored": {
			configVal: types.ListValueMust(
				types.ObjectType{AttrTypes: memberObjectType()},
				[]attr.Value{
					createMemberObject("usr-1", "User", intPtr(1)),
					createMemberObject("usr-2", "User", nil),
					createMemberObject("sch-1", "Schedule", nil),
				},
			),
			expectedErr: false,
		},
		"valid list - unknown list value": {
			configVal:   types.ListUnknown(types.ObjectType{AttrTypes: memberObjectType()}),
			expectedErr: false,
		},
		"invalid list - duplicate positions": {
			configVal: types.ListValueMust(
				types.ObjectType{AttrTypes: memberObjectType()},
				[]attr.Value{
					createMemberObject("usr-1", "User", intPtr(1)),
					createMemberObject("usr-2", "User", intPtr(2)),
					createMemberObject("sch-1", "Schedule", intPtr(1)),
				},
			),
			expectedErr:  true,
			expectedDiag: "Duplicate Attribute Value",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			val := listvalidator.UniqueByAttribute("position")

			req := validator.ListRequest{
				Path:        path.Root("schedule_rotation_members"),
				ConfigValue: test.configVal,
			}
			resp := &validator.ListResponse{}

			val.ValidateList(ctx, req, resp)

			if test.expectedErr && !resp.Diagnostics.HasError() {
				t.Fatalf("expected diagnostics error, got none")
			}

			if !test.expectedErr && resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics error: %s", resp.Diagnostics.Errors()[0].Summary())
			}

			if test.expectedErr && len(resp.Diagnostics.Errors()) > 0 {
				summary := resp.Diagnostics.Errors()[0].Summary()
				if summary != test.expectedDiag {
					t.Errorf("expected diagnostic summary %q, got %q", test.expectedDiag, summary)
				}
			}
		})
	}
}
