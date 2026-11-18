package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func validCSSHexColor() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`),
		"must be a valid action (usually starts with lambda:)",
	)
}

type resourceDiffGetter interface {
	Id() string
	GetOk(key string) (interface{}, bool)
	Get(key string) interface{}
	HasChange(key string) bool
}

func validateTeamAdminIdsInUserIds(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	return validateTeamAdminIdsInUserIdsInternal(ctx, d, meta)
}

func validateTeamAdminIdsInUserIdsInternal(_ context.Context, d resourceDiffGetter, _ interface{}) error {
	adminRaw, adminOk := d.GetOk("admin_ids")
	if !adminOk {
		return nil
	}
	adminList, ok := adminRaw.([]interface{})
	if !ok || len(adminList) == 0 {
		return nil
	}

	userRaw, userOk := d.GetOk("user_ids")
	userSet := make(map[int]bool)
	if userOk {
		if userList, ok := userRaw.([]interface{}); ok {
			for _, u := range userList {
				if id, ok := u.(int); ok {
					userSet[id] = true
				}
			}
		}
	}

	for _, a := range adminList {
		id, ok := a.(int)
		if !ok {
			continue
		}
		if !userSet[id] {
			return fmt.Errorf("admin_ids contains user %d which is not in user_ids — every admin must also be a team member", id)
		}
	}
	return nil
}

func validateUniqueWorkflowTaskPosition(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	return validateUniqueWorkflowTaskPositionInternal(ctx, d, meta)
}

func validateScheduleRotationMemberPositions(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	return validateScheduleRotationMemberPositionsInternal(ctx, d, meta)
}

func validateScheduleRotationMemberPositionsInternal(ctx context.Context, d resourceDiffGetter, meta interface{}) error {
	members, ok := d.GetOk("schedule_rotation_members")
	if !ok {
		return nil
	}

	memberList, ok := members.([]interface{})
	if !ok || len(memberList) == 0 {
		return nil
	}

	seen := make(map[int]bool)
	for i, memberRaw := range memberList {
		member, ok := memberRaw.(map[string]interface{})
		if !ok {
			continue
		}

		posRaw, ok := member["position"]
		if !ok {
			continue
		}
		pos, ok := posRaw.(int)
		if !ok {
			continue
		}

		if pos <= 0 {
			return fmt.Errorf("schedule_rotation_members[%d].position must be greater than 0 (got %d)", i, pos)
		}
		if seen[pos] {
			return fmt.Errorf("schedule_rotation_members has duplicate position %d — positions must be unique and 1-indexed", pos)
		}
		seen[pos] = true
	}

	return nil
}

// validateUniqueWorkflowTaskPositionInternal is the internal implementation that accepts an interface for testing.
// It does not check for positions already in use: the API stores positions as given, and a per-resource check
// can't see other tasks moving in the same plan, so it rejected valid reorders.
func validateUniqueWorkflowTaskPositionInternal(_ context.Context, d resourceDiffGetter, _ interface{}) error {
	position, ok := d.GetOk("position")
	if !ok {
		return nil
	}

	if positionInt, ok := position.(int); ok && positionInt <= 0 {
		return fmt.Errorf("position must be greater than 0")
	}
	return nil
}
