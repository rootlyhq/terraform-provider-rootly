package rootlytypes

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.StringTypable        = (*RFC3339Type)(nil)
	_ basetypes.StringValuable       = (*RFC3339)(nil)
	_ xattr.ValidateableAttribute    = (*RFC3339)(nil)
	_ function.ValidateableParameter = (*RFC3339)(nil)
)

type RFC3339Type struct {
	timetypes.RFC3339Type
}

func (t RFC3339Type) String() string {
	return "rootlytypes.RFC3339Type"
}

func (t RFC3339Type) ValueType(ctx context.Context) attr.Value {
	return RFC3339{}
}

func (t RFC3339Type) Equal(o attr.Type) bool {
	other, ok := o.(RFC3339Type)

	if !ok {
		return false
	}

	return t.RFC3339Type.Equal(other.RFC3339Type)
}

func (t RFC3339Type) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return RFC3339{
		RFC3339: timetypes.RFC3339{
			StringValue: in,
		},
	}, nil
}

func (t RFC3339Type) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.RFC3339Type.ValueFromTerraform(ctx, in)

	if err != nil {
		return nil, err
	}

	rfc3339Value, ok := attrValue.(timetypes.RFC3339)

	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	return RFC3339{
		RFC3339: rfc3339Value,
	}, nil
}

type RFC3339 struct {
	timetypes.RFC3339
}

func (v RFC3339) Type(_ context.Context) attr.Type {
	return RFC3339Type{}
}

func (v RFC3339) Equal(o attr.Value) bool {
	other, ok := o.(RFC3339)

	if !ok {
		return false
	}

	return v.RFC3339.Equal(other.RFC3339)
}

func (v RFC3339) StringSemanticEquals(_ context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(RFC3339)
	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", v)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newValuable),
		)

		return false, diags
	}

	// RFC3339 strings are already validated at this point, ignoring errors
	newRFC3339time, _ := time.Parse(time.RFC3339, newValue.ValueString())
	currentRFC3339time, _ := time.Parse(time.RFC3339, v.ValueString())

	return currentRFC3339time.Equal(newRFC3339time), diags
}

func NewRFC3339ValueMust(value string) RFC3339 {
	return RFC3339{
		RFC3339: timetypes.NewRFC3339ValueMust(value),
	}
}
