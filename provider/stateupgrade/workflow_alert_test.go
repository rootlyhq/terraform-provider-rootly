package stateupgrade

import (
	"context"
	"reflect"
	"testing"
)

func TestUpgradeWorkflowAlertV0ToV1_HappyPath(t *testing.T) {
	v0 := map[string]any{
		"name": "ir-4120-migration-test",
		"trigger_params": []any{
			map[string]any{
				"triggers": []any{"alert_created"},
				"alert_payload_conditions": map[string]any{
					"logic":      "ALL",
					"conditions": `[{"query":"$.commonLabels.namespace","operator":"IS","values":["production"]},{"query":"$.commonLabels.severity","operator":"CONTAINS","values":["critical","high"]}]`,
				},
			},
		},
	}

	got, err := UpgradeWorkflowAlertV0ToV1(context.Background(), v0, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	apc := got["trigger_params"].([]any)[0].(map[string]any)["alert_payload_conditions"]
	list, ok := apc.([]any)
	if !ok {
		t.Fatalf("expected []any, got %T", apc)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 nested block, got %d", len(list))
	}
	block := list[0].(map[string]any)
	if block["logic"] != "ALL" {
		t.Errorf("logic = %v, want ALL", block["logic"])
	}
	conds := block["conditions"].([]any)
	if len(conds) != 2 {
		t.Fatalf("expected 2 conditions, got %d", len(conds))
	}
	c0 := conds[0].(map[string]any)
	if c0["query"] != "$.commonLabels.namespace" || c0["operator"] != "IS" {
		t.Errorf("condition 0 mismatch: %+v", c0)
	}
	if !reflect.DeepEqual(c0["values"], []any{"production"}) {
		t.Errorf("condition 0 values = %v", c0["values"])
	}
	c1 := conds[1].(map[string]any)
	if !reflect.DeepEqual(c1["values"], []any{"critical", "high"}) {
		t.Errorf("condition 1 values = %v", c1["values"])
	}
}

func TestUpgradeWorkflowAlertV0ToV1_UseRegexpPreserved(t *testing.T) {
	v0 := map[string]any{
		"trigger_params": []any{
			map[string]any{
				"alert_payload_conditions": map[string]any{
					"logic":      "ANY",
					"conditions": `[{"query":"$.alertname","operator":"IS","values":["^api-.+"],"use_regexp":true}]`,
				},
			},
		},
	}
	got, err := UpgradeWorkflowAlertV0ToV1(context.Background(), v0, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	c := got["trigger_params"].([]any)[0].(map[string]any)["alert_payload_conditions"].([]any)[0].(map[string]any)["conditions"].([]any)[0].(map[string]any)
	if c["use_regexp"] != true {
		t.Errorf("use_regexp = %v, want true", c["use_regexp"])
	}
}

func TestUpgradeWorkflowAlertV0ToV1_EmptyMapRemoved(t *testing.T) {
	v0 := map[string]any{
		"trigger_params": []any{
			map[string]any{
				"triggers":                 []any{"alert_created"},
				"alert_payload_conditions": map[string]any{},
			},
		},
	}
	got, err := UpgradeWorkflowAlertV0ToV1(context.Background(), v0, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tp := got["trigger_params"].([]any)[0].(map[string]any)
	if _, present := tp["alert_payload_conditions"]; present {
		t.Errorf("expected empty alert_payload_conditions to be removed")
	}
}

func TestUpgradeWorkflowAlertV0ToV1_MissingTriggerParams(t *testing.T) {
	v0 := map[string]any{"name": "no-tp"}
	got, err := UpgradeWorkflowAlertV0ToV1(context.Background(), v0, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, v0) {
		t.Errorf("state was modified unexpectedly: %+v", got)
	}
}

func TestUpgradeWorkflowAlertV0ToV1_MalformedConditionsJSON(t *testing.T) {
	v0 := map[string]any{
		"trigger_params": []any{
			map[string]any{
				"alert_payload_conditions": map[string]any{
					"logic":      "ALL",
					"conditions": `{not valid json`,
				},
			},
		},
	}
	got, err := UpgradeWorkflowAlertV0ToV1(context.Background(), v0, nil)
	if err != nil {
		t.Fatalf("upgrader should not fail on malformed JSON, got: %v", err)
	}
	block := got["trigger_params"].([]any)[0].(map[string]any)["alert_payload_conditions"].([]any)[0].(map[string]any)
	if block["logic"] != "ALL" {
		t.Errorf("logic = %v, want ALL", block["logic"])
	}
	if block["conditions"] != nil {
		if conds, ok := block["conditions"].([]any); !ok || len(conds) != 0 {
			t.Errorf("expected empty conditions, got %+v", block["conditions"])
		}
	}
}

func TestUpgradeWorkflowAlertV0ToV1_AlreadyV1(t *testing.T) {
	v1 := map[string]any{
		"trigger_params": []any{
			map[string]any{
				"alert_payload_conditions": []any{
					map[string]any{
						"logic":      "ALL",
						"conditions": []any{map[string]any{"query": "$.x", "operator": "IS"}},
					},
				},
			},
		},
	}
	got, err := UpgradeWorkflowAlertV0ToV1(context.Background(), v1, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, v1) {
		t.Errorf("already-V1 state was modified: %+v", got)
	}
}

func TestWorkflowAlertV0_ImpliedTypeDoesNotPanic(t *testing.T) {
	_ = WorkflowAlertV0().CoreConfigSchema().ImpliedType()
}
