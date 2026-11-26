package provider

import (
	"context"
	"testing"
)

func TestValidPermissionAction(t *testing.T) {
	validValues := []string{
		"#112233",
		"#123",
		"#000233",
		"#023",
	}
	for _, v := range validValues {
		_, errors := validCSSHexColor()(v, "action")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid CSS Hex color: %q", v, errors)
		}
	}

	invalidNames := []string{
		"invalid",
		"#abcd",
		"#-12",
	}
	for _, v := range invalidNames {
		_, errors := validCSSHexColor()(v, "action")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid CSS Hex color", v)
		}
	}
}

func TestValidateUniqueWorkflowTaskPosition_NonPositivePosition(t *testing.T) {
	d := &ResourceDiffMock{
		Id_:         "my-task-1",
		Values:      map[string]interface{}{"position": -1, "workflow_id": "wf1"},
		ChangedKeys: map[string]bool{"position": true},
	}
	if err := validateUniqueWorkflowTaskPositionInternal(context.Background(), d, nil); err == nil {
		t.Fatalf("expected an error for a negative position, but got nil")
	}
}

type ResourceDiffMock struct {
	Id_         string
	Values      map[string]interface{}
	ChangedKeys map[string]bool
}

func (r *ResourceDiffMock) Id() string { return r.Id_ }
func (r *ResourceDiffMock) GetOk(key string) (interface{}, bool) {
	v, ok := r.Values[key]
	return v, ok
}
func (r *ResourceDiffMock) Get(key string) interface{} { return r.Values[key] }
func (r *ResourceDiffMock) HasChange(key string) bool  { return r.ChangedKeys[key] }
func (r *ResourceDiffMock) GetChange(key string) (interface{}, interface{}) {
	return nil, r.Values[key]
}

func TestValidateScheduleRotationMemberPositions_Valid(t *testing.T) {
	d := &ResourceDiffMock{
		Values: map[string]interface{}{
			"schedule_rotation_members": []interface{}{
				map[string]interface{}{"position": 1, "user_id": "user-1"},
				map[string]interface{}{"position": 2, "user_id": "user-2"},
				map[string]interface{}{"position": 3, "user_id": "user-3"},
			},
		},
	}
	err := validateScheduleRotationMemberPositionsInternal(context.Background(), d, nil)
	if err != nil {
		t.Fatalf("expected no error for valid sequential positions, got: %v", err)
	}
}

func TestValidateScheduleRotationMemberPositions_DuplicatePosition(t *testing.T) {
	d := &ResourceDiffMock{
		Values: map[string]interface{}{
			"schedule_rotation_members": []interface{}{
				map[string]interface{}{"position": 1, "user_id": "user-1"},
				map[string]interface{}{"position": 1, "user_id": "user-2"},
			},
		},
	}
	err := validateScheduleRotationMemberPositionsInternal(context.Background(), d, nil)
	if err == nil {
		t.Fatalf("expected error for duplicate position, got nil")
	}
}

func TestValidateScheduleRotationMemberPositions_ZeroPosition(t *testing.T) {
	d := &ResourceDiffMock{
		Values: map[string]interface{}{
			"schedule_rotation_members": []interface{}{
				map[string]interface{}{"position": 0, "user_id": "user-1"},
				map[string]interface{}{"position": 1, "user_id": "user-2"},
			},
		},
	}
	err := validateScheduleRotationMemberPositionsInternal(context.Background(), d, nil)
	if err == nil {
		t.Fatalf("expected error for position 0, got nil")
	}
}

func TestValidateScheduleRotationMemberPositions_NegativePosition(t *testing.T) {
	d := &ResourceDiffMock{
		Values: map[string]interface{}{
			"schedule_rotation_members": []interface{}{
				map[string]interface{}{"position": -1, "user_id": "user-1"},
			},
		},
	}
	err := validateScheduleRotationMemberPositionsInternal(context.Background(), d, nil)
	if err == nil {
		t.Fatalf("expected error for negative position, got nil")
	}
}

func TestValidateScheduleRotationMemberPositions_NoMembers(t *testing.T) {
	d := &ResourceDiffMock{
		Values: map[string]interface{}{},
	}
	err := validateScheduleRotationMemberPositionsInternal(context.Background(), d, nil)
	if err != nil {
		t.Fatalf("expected no error when no members set, got: %v", err)
	}
}

func TestValidateTeamAdminIdsInUserIds_Valid(t *testing.T) {
	d := &ResourceDiffMock{
		Values: map[string]interface{}{
			"user_ids":  []interface{}{1, 2, 3},
			"admin_ids": []interface{}{1, 3},
		},
	}
	err := validateTeamAdminIdsInUserIdsInternal(context.Background(), d, nil)
	if err != nil {
		t.Fatalf("expected no error when all admins are in user_ids: %v", err)
	}
}

func TestValidateTeamAdminIdsInUserIds_AdminNotInUsers(t *testing.T) {
	d := &ResourceDiffMock{
		Values: map[string]interface{}{
			"user_ids":  []interface{}{1, 2},
			"admin_ids": []interface{}{1, 99},
		},
	}
	err := validateTeamAdminIdsInUserIdsInternal(context.Background(), d, nil)
	if err == nil {
		t.Fatalf("expected error when admin 99 is not in user_ids, got nil")
	}
}

func TestValidateTeamAdminIdsInUserIds_NoAdmins(t *testing.T) {
	d := &ResourceDiffMock{
		Values: map[string]interface{}{
			"user_ids": []interface{}{1, 2},
		},
	}
	err := validateTeamAdminIdsInUserIdsInternal(context.Background(), d, nil)
	if err != nil {
		t.Fatalf("expected no error when no admin_ids set: %v", err)
	}
}

func TestValidateTeamAdminIdsInUserIds_NoUsers(t *testing.T) {
	d := &ResourceDiffMock{
		Values: map[string]interface{}{
			"admin_ids": []interface{}{1},
		},
	}
	err := validateTeamAdminIdsInUserIdsInternal(context.Background(), d, nil)
	if err == nil {
		t.Fatalf("expected error when admin set but no user_ids, got nil")
	}
}

// A valid reorder (task "a" moved from 1 to 5, the rest shifted up)
// should be accepted.
func TestValidateUniqueWorkflowTaskPosition_ReorderDeadlock(t *testing.T) {
	// Live positions are a=1 b=2 c=3 d=4 e=5; each planned position collides with a neighbour that is also moving.
	planned := map[string]int{"a": 5, "b": 1, "c": 2, "d": 3, "e": 4}

	for id, position := range planned {
		d := &ResourceDiffMock{
			Id_:         id,
			Values:      map[string]interface{}{"position": position, "workflow_id": "wf1"},
			ChangedKeys: map[string]bool{"position": true},
		}
		if err := validateUniqueWorkflowTaskPositionInternal(context.Background(), d, nil); err != nil {
			t.Errorf("task %s -> %d: %v", id, position, err)
		}
	}
}
