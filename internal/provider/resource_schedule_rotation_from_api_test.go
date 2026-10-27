package provider

import (
	"context"
	"testing"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	"github.com/rootlyhq/terraform-provider-rootly/v5/internal/apiclient"
	"github.com/stretchr/testify/require"
)

// After import, state has no members yet; FromApi must still read them from the API.
func TestScheduleRotationFromApiReadsMembersAfterImport(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	m := ScheduleRotationResourceModel{
		ScheduleRotationMembers: supertypes.NewSetNestedObjectValueOfNull[ScheduleRotationResourceScheduleRotationMemberModel](ctx),
	}

	diags := m.FromApi(ctx, apiclient.ScheduleRotation{
		ID: "rotation-id",
		ScheduleRotationMembers: []apiclient.ScheduleRotationMember{
			{MemberType: "User", MemberID: "1", Position: 1},
			{MemberType: "User", MemberID: "2", Position: 2},
		},
	})
	require.False(t, diags.HasError(), diags)

	members, diags := m.ScheduleRotationMembers.Get(ctx)
	require.False(t, diags.HasError(), diags)
	require.Len(t, members, 2)
}
