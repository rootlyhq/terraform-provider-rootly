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
