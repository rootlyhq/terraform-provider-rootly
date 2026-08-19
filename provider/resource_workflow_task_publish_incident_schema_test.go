package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceWorkflowTaskPublishIncidentAcceptsComponentSelection(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, resourceWorkflowTaskPublishIncident().Schema, nil)
	taskParams := []interface{}{
		map[string]interface{}{
			"selected_component_keys": []interface{}{"Service:123"},
			"selected_component_statuses": map[string]interface{}{
				"Service:123": "major_outage",
			},
		},
	}

	if err := resourceData.Set("task_params", taskParams); err != nil {
		t.Fatalf("set task_params returned an error: %v", err)
	}
}
