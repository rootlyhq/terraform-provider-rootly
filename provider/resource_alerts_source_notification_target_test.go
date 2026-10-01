package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v5/client"
)

func TestAlertsSourceNotificationTargetOnWire(t *testing.T) {
	var sent map[string]interface{}
	returned := map[string]interface{}{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.api+json")
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			var body struct {
				Data struct {
					Attributes struct {
						SourceableAttributes map[string]interface{} `json:"sourceable_attributes"`
					} `json:"attributes"`
				} `json:"data"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decoding %s request body: %v", r.Method, err)
			}
			sent = body.Data.Attributes.SourceableAttributes
			returned = map[string]interface{}{
				"notification_target_type": sent["notification_target_type"],
				"notification_target_id":   sent["notification_target_id"],
			}
			if sent["notification_target_id"] == "" {
				returned = map[string]interface{}{"notification_target_type": nil, "notification_target_id": nil}
			}
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"id": "1", "type": "alert_sources",
				"attributes": map[string]interface{}{"source_type": "email", "sourceable_attributes": returned},
			},
		})
	}))
	t.Cleanup(server.Close)

	api, err := client.NewClient(server.URL, "test-token", "alerts-source-test")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	resource := resourceAlertsSource()

	data := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"name":        "email",
		"source_type": "email",
		"sourceable_attributes": []interface{}{map[string]interface{}{
			"notification_target_type": "Group",
			"notification_target_id":   "group-1",
		}},
	})
	if diags := resourceAlertsSourceCreate(ctx, data, api); diags.HasError() {
		t.Fatal(diags)
	}
	if sent["notification_target_type"] != "Group" || sent["notification_target_id"] != "group-1" {
		t.Fatalf("create sent %#v", sent)
	}
	if got := data.Get("sourceable_attributes.0.notification_target_id"); got != "group-1" {
		t.Fatalf("expected notification_target_id group-1 in state, got %#v", got)
	}
	if got := data.Get("sourceable_attributes.0.notification_target_type"); got != "Group" {
		t.Fatalf("expected notification_target_type Group in state, got %#v", got)
	}

	if err := data.Set("sourceable_attributes", []interface{}{map[string]interface{}{
		"notification_target_type": "",
		"notification_target_id":   "",
	}}); err != nil {
		t.Fatal(err)
	}
	if diags := resourceAlertsSourceUpdate(ctx, data, api); diags.HasError() {
		t.Fatal(diags)
	}
	if v, ok := sent["notification_target_id"]; !ok || v != "" {
		t.Fatalf("clearing should send an empty notification_target_id, sent %#v", sent)
	}
	if got := data.Get("sourceable_attributes.0.notification_target_id"); got != "" {
		t.Fatalf("expected cleared notification_target_id in state, got %#v", got)
	}
}
