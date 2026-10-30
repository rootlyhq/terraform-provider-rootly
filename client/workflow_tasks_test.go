package client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// The parent workflow goes in the URL path. The API rejects workflow_id in the
// request body with a 400, so create and update must not send it (#525).
func TestWorkflowTaskRequestsOmitWorkflowIdAttribute(t *testing.T) {
	requests := map[string]map[string]interface{}{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("%s %s: reading body: %v", r.Method, r.URL.Path, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var body struct {
			Data struct {
				Attributes map[string]interface{} `json:"attributes"`
			} `json:"data"`
		}
		if err := json.Unmarshal(raw, &body); err != nil {
			t.Errorf("%s %s: invalid body %q: %v", r.Method, r.URL.Path, raw, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		requests[r.Method+" "+r.URL.Path] = body.Data.Attributes

		w.Header().Set("Content-Type", "application/vnd.api+json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{
			"id": "task-id", "type": "workflow_tasks", "attributes": map[string]interface{}{
				"workflow_id": "workflow-id", "name": "Timeline", "position": 1,
				"task_params": map[string]interface{}{"task_type": "add_to_timeline", "event": "Test"},
			},
		}})
	}))
	defer server.Close()

	c, err := NewClient(server.URL+"/api", "test-token", "workflow-task-test")
	if err != nil {
		t.Fatal(err)
	}

	task := &WorkflowTask{
		WorkflowId: "workflow-id",
		Name:       "Timeline",
		Position:   1,
		TaskParams: map[string]interface{}{"task_type": "add_to_timeline", "event": "Test"},
	}

	created, err := c.CreateWorkflowTask(task)
	if err != nil {
		t.Fatal(err)
	}
	if created.WorkflowId != "workflow-id" {
		t.Errorf("expected workflow_id to be read from the create response, got %q", created.WorkflowId)
	}
	if task.WorkflowId != "workflow-id" {
		t.Errorf("CreateWorkflowTask must not modify its argument, workflow_id is now %q", task.WorkflowId)
	}

	updated, err := c.UpdateWorkflowTask("task-id", task)
	if err != nil {
		t.Fatal(err)
	}
	if updated.WorkflowId != "workflow-id" {
		t.Errorf("expected workflow_id to be read from the update response, got %q", updated.WorkflowId)
	}

	for _, key := range []string{
		http.MethodPost + " /api/v1/workflows/workflow-id/workflow_tasks",
		http.MethodPut + " /api/v1/workflow_tasks/task-id",
	} {
		attrs, ok := requests[key]
		if !ok {
			t.Errorf("expected request %s, got %v", key, requests)
			continue
		}
		if _, present := attrs["workflow_id"]; present {
			t.Errorf("%s: workflow_id must not be sent as an attribute, got %v", key, attrs)
		}
		if attrs["name"] != "Timeline" {
			t.Errorf("%s: expected name attribute to be sent, got %v", key, attrs)
		}
	}
}
