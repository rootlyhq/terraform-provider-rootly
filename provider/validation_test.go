package provider

import (
	"context"
	"testing"

	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v2/schema"
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

type dummyClient struct {
	tasks []*struct {
		ID       string
		Name     string
		Position int
	}
}

func (c *dummyClient) ListWorkflowTasks(workflowId string, params *rootlygo.ListWorkflowTasksParams) ([]interface{}, error) {
	var res []interface{}
	for _, t := range c.tasks {
		res = append(res, &client.WorkflowTask{
			ID:       t.ID,
			Name:     t.Name,
			Position: t.Position,
		})
	}
	return res, nil
}

func TestValidateUniqueWorkflowTaskPosition_RunValidationDirectly(t *testing.T) {
	// Create a test ResourceDiffMock
	d := &ResourceDiffMock{
		Id_:         "my-task-2",
		Values:      map[string]interface{}{"position": 1, "workflow_id": "wf1"},
		ChangedKeys: map[string]bool{"position": true, "workflow_id": true},
	}
	c := &dummyClient{
		tasks: []*struct {
			ID       string
			Name     string
			Position int
		}{
			{ID: "my-task-1", Name: "First Task", Position: 1}, // Conflict
		},
	}
	err := validateUniqueWorkflowTaskPositionInternal(context.Background(), d, c)
	if err == nil {
		t.Fatalf("expected an error due to duplicate position, but got nil")
	}
	d = &ResourceDiffMock{
		Id_:         "my-task-3",
		Values:      map[string]interface{}{"position": 2, "workflow_id": "wf1"},
		ChangedKeys: map[string]bool{"position": true, "workflow_id": true},
	}
	err = validateUniqueWorkflowTaskPositionInternal(context.Background(), d, c)
	if err != nil {
		t.Fatalf("expected no error when position is unique: %v", err)
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
