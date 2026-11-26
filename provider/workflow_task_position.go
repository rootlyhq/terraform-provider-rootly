package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v5/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v5/schema"
)

// workflowTaskLister is an interface for listing workflow tasks (for testing)
type workflowTaskLister interface {
	ListWorkflowTasks(workflowId string, params *rootlygo.ListWorkflowTasksParams) ([]interface{}, error)
}

// readWorkflowTaskWithPositionWarning wraps a workflow task Read so that a refresh warns when another task in the
// same workflow has the same position. Create and Update call the inner Read directly, so tasks that share a
// position only for the moment while a reorder is being applied don't warn.
func readWorkflowTaskWithPositionWarning(read schema.ReadContextFunc) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		diags := read(ctx, d, meta)
		if diags.HasError() || d.Id() == "" {
			return diags
		}

		lister, ok := meta.(workflowTaskLister)
		if !ok {
			return diags
		}

		workflowId, _ := d.Get("workflow_id").(string)
		position, _ := d.Get("position").(int)
		return append(diags, workflowTaskPositionWarnings(ctx, lister, workflowId, d.Id(), position)...)
	}
}

func workflowTaskPositionWarnings(ctx context.Context, lister workflowTaskLister, workflowId string, taskId string, position int) diag.Diagnostics {
	if workflowId == "" || position <= 0 {
		return nil
	}

	tasks, err := lister.ListWorkflowTasks(workflowId, &rootlygo.ListWorkflowTasksParams{})
	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("Could not check positions of workflow %s: %v", workflowId, err))
		return nil
	}

	var diags diag.Diagnostics
	for _, taskInterface := range tasks {
		task, ok := taskInterface.(*client.WorkflowTask)
		if !ok || task.ID == taskId || task.Position != position {
			continue
		}

		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Warning,
			Summary:       fmt.Sprintf("Workflow task position %d is also used by workflow task %s", position, task.ID),
			Detail:        fmt.Sprintf("Workflow task %s (name: %s) has the same position. Give each task in a workflow a unique position, or a workflow run may skip one of them.", task.ID, task.Name),
			AttributePath: cty.GetAttrPath("position"),
		})
	}
	return diags
}
