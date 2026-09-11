package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v5/client"
)

func TestWorkflowTaskSlackCanvasClearsWorkspaceOnWire(t *testing.T) {
	for _, action := range []struct {
		name     string
		resource func() *schema.Resource
	}{
		{"create_slack_canvas", resourceWorkflowTaskCreateSlackCanvas},
		{"update_slack_canvas", resourceWorkflowTaskUpdateSlackCanvas},
	} {
		t.Run(action.name, func(t *testing.T) {
			updates := 0
			params := map[string]interface{}{"task_type": action.name, "content": "# Updated {{ incident.title }}", "title": "Incident report"}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/v1/workflow_tasks/task-id" {
					t.Errorf("unexpected request path %s", r.URL.Path)
					w.WriteHeader(404)
					return
				}
				if r.Method == http.MethodPut {
					updates++
					var body struct {
						Data struct {
							Attributes struct {
								TaskParams map[string]interface{} `json:"task_params"`
							} `json:"attributes"`
						} `json:"data"`
					}
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						t.Error(err)
						w.WriteHeader(400)
						return
					}
					params = body.Data.Attributes.TaskParams
					channel, ok := params["channel"].(map[string]interface{})
					if !ok {
						t.Error("channel was not a JSON object")
						w.WriteHeader(400)
						return
					}
					workspace, present := channel["workspace"]
					if !present || workspace != nil {
						t.Errorf("PUT must contain workspace:null, got %#v", channel)
					}
					delete(channel, "workspace")
				} else if r.Method != http.MethodGet {
					t.Errorf("unexpected method %s", r.Method)
				}
				w.Header().Set("Content-Type", "application/vnd.api+json")
				json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{
					"id": "task-id", "type": "workflow_tasks", "attributes": map[string]interface{}{
						"workflow_id": "workflow-id", "name": "Canvas", "position": 1, "enabled": false, "skip_on_failure": false, "task_params": params,
					},
				}})
			}))
			defer server.Close()
			api, err := client.NewClient(server.URL+"/api", "test-token", "canvas-test")
			if err != nil {
				t.Fatal(err)
			}
			resource := action.resource()
			values := map[string]interface{}{
				"task_type": action.name, "content": "# Updated {{ incident.title }}",
				"channel": []interface{}{map[string]interface{}{"id": "C123", "name": "incidents"}},
			}
			if action.name == "create_slack_canvas" {
				values["title"] = "Incident report"
			}
			data := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
				"workflow_id": "workflow-id", "name": "Canvas", "enabled": false, "position": 1,
				"task_params": []interface{}{values},
			})
			data.SetId("task-id")
			if diagnostics := resource.UpdateContext(context.Background(), data, api); diagnostics.HasError() {
				t.Fatal(diagnostics)
			}
			if updates != 1 {
				t.Fatalf("expected one PUT, got %d", updates)
			}
			if got := data.Get("task_params.0.channel.0.workspace.#"); fmt.Sprint(got) != "0" {
				t.Fatalf("workspace remained in state: %v", got)
			}
		})
	}
}
