package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAlertRoutingRuleAiSreProviderAcceptsEmptyStringToDisable(t *testing.T) {
	validate := resourceAlertRoutingRule().Schema["ai_sre_provider"].ValidateFunc
	for _, value := range []string{"", "rootly", "devin"} {
		if _, errs := validate(value, "ai_sre_provider"); len(errs) > 0 {
			t.Errorf("expected %q to be accepted, got %v", value, errs)
		}
	}
	if _, errs := validate("other", "ai_sre_provider"); len(errs) == 0 {
		t.Error("expected an unknown provider to be rejected")
	}
}

func TestFormSetResourceTypeForcesReplacement(t *testing.T) {
	if !resourceFormSet().Schema["resource_type"].ForceNew {
		t.Error("expected resource_type to force a new form set")
	}
}

func TestPublishIncidentSyncIncidentComponentsIsComputed(t *testing.T) {
	taskParams := resourceWorkflowTaskPublishIncident().Schema["task_params"].Elem.(*schema.Resource).Schema
	sync := taskParams["sync_incident_components"]
	if !sync.Optional || !sync.Computed {
		t.Errorf("expected sync_incident_components to be optional and computed, got optional=%v computed=%v", sync.Optional, sync.Computed)
	}
}

func TestNotionPageTaskExportSectionsDefaultToIncluded(t *testing.T) {
	resources := map[string]*schema.Resource{
		"create_notion_page": resourceWorkflowTaskCreateNotionPage(),
		"update_notion_page": resourceWorkflowTaskUpdateNotionPage(),
	}
	fields := []string{"include_timeline", "include_follow_ups"}

	for name, resource := range resources {
		taskParamsSchema := resource.Schema["task_params"].Elem.(*schema.Resource).Schema
		for _, field := range fields {
			property := taskParamsSchema[field]
			if property == nil || !property.Optional || property.Default != true {
				t.Errorf("%s: expected %s to be optional with default true, got %#v", name, field, property)
			}
		}

		// Task params are sent as the whole map, so an unset field must be
		// sent as true rather than false.
		unset := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
			"workflow_id": "workflow",
			"task_params": []interface{}{map[string]interface{}{"title": "Retrospective"}},
		})
		params := unset.Get("task_params").([]interface{})[0].(map[string]interface{})
		for _, field := range fields {
			if params[field] != true {
				t.Errorf("%s: expected unset %s to be sent as true, got %#v", name, field, params[field])
			}
		}

		disabled := map[string]interface{}{"title": "Retrospective"}
		for _, field := range fields {
			disabled[field] = false
		}
		set := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
			"workflow_id": "workflow",
			"task_params": []interface{}{disabled},
		})
		params = set.Get("task_params").([]interface{})[0].(map[string]interface{})
		for _, field := range fields {
			if params[field] != false {
				t.Errorf("%s: expected %s = false to be sent as false, got %#v", name, field, params[field])
			}
		}
	}
}
