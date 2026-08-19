package provider

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func validatePublishIncidentWorkflowTaskDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	if err := validateUniqueWorkflowTaskPosition(ctx, d, meta); err != nil {
		return err
	}

	taskParams, ok := d.GetOk("task_params")
	if !ok {
		return nil
	}
	paramsList, ok := taskParams.([]interface{})
	if !ok || len(paramsList) == 0 || paramsList[0] == nil {
		return nil
	}
	params, ok := paramsList[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return validatePublishIncidentComponentSelection(params)
}

func validatePublishIncidentComponentSelection(params map[string]interface{}) error {
	selectedKeys := make(map[string]struct{})
	if keys, ok := params["selected_component_keys"].([]interface{}); ok {
		for _, key := range keys {
			selectedKeys[fmt.Sprint(key)] = struct{}{}
		}
	}

	statuses, ok := params["selected_component_statuses"].(map[string]interface{})
	if !ok {
		return nil
	}
	unknownKeys := make([]string, 0)
	for key := range statuses {
		if _, selected := selectedKeys[key]; !selected {
			unknownKeys = append(unknownKeys, key)
		}
	}
	if len(unknownKeys) == 0 {
		return nil
	}

	sort.Strings(unknownKeys)
	return fmt.Errorf(
		"task_params.0.selected_component_statuses contains keys absent from selected_component_keys: %v",
		unknownKeys,
	)
}
