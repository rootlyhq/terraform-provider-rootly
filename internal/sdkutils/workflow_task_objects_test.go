package sdkutils

import (
	"reflect"
	"testing"
)

func workflowTaskObjectTestFields() WorkflowTaskObjectFields {
	return WorkflowTaskObjectFields{
		"channel": {
			"workspace": {},
		},
	}
}

func TestWorkflowTaskObjectsRoundTrip(t *testing.T) {
	fields := workflowTaskObjectTestFields()
	api := map[string]interface{}{
		"channel": map[string]interface{}{
			"id": "C123", "name": "incidents",
			"workspace": map[string]interface{}{"id": "T123", "name": "Engineering"},
		},
		"content": "# {{ incident.title }}",
		"flat":    map[string]interface{}{"id": "unchanged", "name": "flat map"},
	}
	state, err := WorkflowTaskObjectsToBlocks(api, fields)
	if err != nil {
		t.Fatal(err)
	}
	actual, err := WorkflowTaskBlocksToObjects(state, fields)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, api) {
		t.Fatalf("round trip changed parameters: %#v", actual)
	}
	if _, ok := api["channel"].(map[string]interface{}); !ok {
		t.Fatal("input was mutated")
	}
	if _, ok := state["flat"].(map[string]interface{}); !ok {
		t.Fatal("flat map was converted to a block")
	}
}

func TestWorkflowTaskObjectsClearWorkspace(t *testing.T) {
	fields := workflowTaskObjectTestFields()
	for _, value := range []interface{}{nil, []interface{}{}} {
		state := map[string]interface{}{"channel": []interface{}{map[string]interface{}{
			"id": "C123", "name": "incidents", "workspace": value,
		}}}
		api, err := WorkflowTaskBlocksToObjects(state, fields)
		if err != nil {
			t.Fatal(err)
		}
		channel := api["channel"].(map[string]interface{})
		workspace, present := channel["workspace"]
		if !present || workspace != nil {
			t.Fatalf("clear must preserve explicit null: %#v", channel)
		}
		delete(channel, "workspace")
		reloaded, err := WorkflowTaskObjectsToBlocks(api, fields)
		if err != nil {
			t.Fatal(err)
		}
		block := reloaded["channel"].([]interface{})[0].(map[string]interface{})
		if len(block["workspace"].([]interface{})) != 0 {
			t.Fatal("omitted workspace must restore an empty block")
		}
	}
}

func TestWorkflowTaskObjectsRejectInvalidShapes(t *testing.T) {
	fields := workflowTaskObjectTestFields()
	for _, value := range []interface{}{"invalid", []interface{}{"invalid"}, []interface{}{map[string]interface{}{}, map[string]interface{}{}}} {
		if _, err := WorkflowTaskBlocksToObjects(map[string]interface{}{"channel": value}, fields); err == nil {
			t.Fatalf("accepted invalid block shape: %T", value)
		}
	}
	if _, err := WorkflowTaskObjectsToBlocks(map[string]interface{}{"channel": "invalid"}, fields); err == nil {
		t.Fatal("accepted invalid API object")
	}
}
