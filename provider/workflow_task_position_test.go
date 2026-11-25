package provider

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/rootlyhq/terraform-provider-rootly/v5/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v5/schema"
)

type fakeWorkflowTaskLister struct {
	tasks []*client.WorkflowTask
	err   error
}

func (l *fakeWorkflowTaskLister) ListWorkflowTasks(workflowId string, params *rootlygo.ListWorkflowTasksParams) ([]interface{}, error) {
	if l.err != nil {
		return nil, l.err
	}
	var res []interface{}
	for _, t := range l.tasks {
		res = append(res, t)
	}
	return res, nil
}

func TestWorkflowTaskPositionWarnings(t *testing.T) {
	lister := &fakeWorkflowTaskLister{
		tasks: []*client.WorkflowTask{
			{ID: "a", Name: "Task A", Position: 1},
			{ID: "b", Name: "Task B", Position: 2},
			{ID: "c", Name: "Task C", Position: 2},
		},
	}

	t.Run("unique position", func(t *testing.T) {
		if diags := workflowTaskPositionWarnings(context.Background(), lister, "wf1", "a", 1); len(diags) != 0 {
			t.Fatalf("expected no diagnostics, got %v", diags)
		}
	})

	t.Run("shared position", func(t *testing.T) {
		diags := workflowTaskPositionWarnings(context.Background(), lister, "wf1", "b", 2)
		if len(diags) != 1 {
			t.Fatalf("expected 1 diagnostic, got %v", diags)
		}
		if diags[0].Severity != diag.Warning {
			t.Fatalf("expected a warning, got severity %v", diags[0].Severity)
		}
		if !strings.Contains(diags[0].Summary, "workflow task c") {
			t.Fatalf("expected the warning to name task c, got %q", diags[0].Summary)
		}
	})

	t.Run("list error", func(t *testing.T) {
		failing := &fakeWorkflowTaskLister{err: errors.New("boom")}
		if diags := workflowTaskPositionWarnings(context.Background(), failing, "wf1", "b", 2); len(diags) != 0 {
			t.Fatalf("expected no diagnostics when listing fails, got %v", diags)
		}
	})
}
