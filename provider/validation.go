package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	rootlygo_ "github.com/rootlyhq/rootly-go"
	"github.com/rootlyhq/terraform-provider-rootly/v5/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v5/schema"
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

// workflowTaskLister is an interface for listing workflow tasks (for testing)
type workflowTaskLister interface {
	ListWorkflowTasks(workflowId rootlygo_.ID, params *rootlygo.ListWorkflowTasksParams) ([]interface{}, error)
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

// validateUniqueWorkflowTaskPositionInternal is the internal implementation that accepts an interface for testing
func validateUniqueWorkflowTaskPositionInternal(ctx context.Context, d resourceDiffGetter, meta interface{}) error {
	if d.Id() == "" && !d.HasChange("workflow_id") {
		return nil
	}

	position, positionExists := d.GetOk("position")
	if !positionExists {
		return nil
	}

	positionInt, ok := position.(int)
	if !ok {
		return nil
	}

	if positionInt <= 0 {
		return fmt.Errorf("position must be greater than 0")
	}
	workflowIdRaw := d.Get("workflow_id")
	workflowId, ok := workflowIdRaw.(string)
	if !ok {
		return fmt.Errorf("workflow_id is required")
	}

	currentTaskId := d.Id()

	// Try to get client - either real client or test mock
	var lister workflowTaskLister
	if c, ok := meta.(*client.Client); ok {
		lister = c
	} else if mockLister, ok := meta.(workflowTaskLister); ok {
		lister = mockLister
	} else {
		return nil
	}

	tasks, err := lister.ListWorkflowTasks(rootlygo_.ID(workflowId), &rootlygo.ListWorkflowTasksParams{})
	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("Could not validate unique position for workflow %s: %v", workflowId, err))
		return nil
	}

	for _, taskInterface := range tasks {
		task, ok := taskInterface.(*client.WorkflowTask)
		if !ok {
			continue
		}

		if currentTaskId != "" && task.ID == currentTaskId {
			continue
		}

		if task.Position == positionInt {
			return fmt.Errorf(
				"position %d is already in use by workflow task %s (name: %s). Each workflow task must have a unique position",
				positionInt,
				task.ID,
				task.Name,
			)
		}
	}
	return nil
}
