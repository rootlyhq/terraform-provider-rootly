package provider

import (
	"reflect"
	"testing"
)

func TestNormalizePageRootlyOnCallRespondersTaskParamsDropsEmptyTargets(t *testing.T) {
	input := map[string]interface{}{
		"task_type":                "page_rootly_on_call_responders",
		"alert_urgency_id":         "urgency-id",
		"summary":                  "Incident: {{ incident.title }}",
		"description":              "page the team",
		"escalation_note":          "",
		"escalation_policy_target": map[string]interface{}{},
		"service_target":           map[string]string{},
		"user_target":              nil,
		"group_target": map[string]interface{}{
			"id":   "new-group-id",
			"name": "New Team",
		},
		"functionality_target": map[interface{}]interface{}{},
	}

	got := normalizePageRootlyOnCallRespondersTaskParams(input)
	want := map[string]interface{}{
		"task_type":        "page_rootly_on_call_responders",
		"alert_urgency_id": "urgency-id",
		"summary":          "Incident: {{ incident.title }}",
		"description":      "page the team",
		"escalation_note":  "",
		"group_target": map[string]interface{}{
			"id":   "new-group-id",
			"name": "New Team",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected task params\n got: %#v\nwant: %#v", got, want)
	}
}

func TestNormalizePageRootlyOnCallRespondersTaskParamsKeepsAllAttachedTeamsTarget(t *testing.T) {
	input := map[string]interface{}{
		"group_target": map[string]interface{}{
			"id":   `{{ incident.raw_groups | map: "id" | join: "," }}`,
			"name": "All Attached Teams",
		},
	}

	got := normalizePageRootlyOnCallRespondersTaskParams(input)
	if !reflect.DeepEqual(got, input) {
		t.Fatalf("unexpected task params\n got: %#v\nwant: %#v", got, input)
	}
}
