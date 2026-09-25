package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/rootlyhq/terraform-provider-rootly/v5/client"
)

type retriggerCase struct {
	name       string
	stateValue string
	configSet  bool
	configVal  int64
	changedKey string
	want       interface{}
}

func retriggerCases() []retriggerCase {
	return []retriggerCase{
		// API returned null, Read flattened it to 0 in state, config still
		// omits the attribute: must send null, not the state's zero value.
		{"update state 0, config absent", "0", false, 0, "name", nil},
		// Clearing a configured override: config omits the attribute again.
		{"update state 30, config absent", "30", false, 0, "name", nil},
		{"update config -1", "30", true, -1, "name", int64(-1)},
		{"update config 30", "0", true, 30, "name", int64(30)},
	}
}

func retriggerRawConfig(set bool, val int64, changedKey string) cty.Value {
	attrs := map[string]cty.Value{
		"retrigger_timeout_minutes": cty.NullVal(cty.Number),
	}
	if set {
		attrs["retrigger_timeout_minutes"] = cty.NumberIntVal(val)
	}
	return cty.ObjectVal(attrs)
}

func retriggerUpdateData(t *testing.T, s *schema.Resource, extraState map[string]string, tc retriggerCase) *schema.ResourceData {
	t.Helper()
	attrs := map[string]string{
		"id":                        "1",
		"retrigger_timeout_minutes": tc.stateValue,
	}
	for k, v := range extraState {
		attrs[k] = v
	}
	diff := &terraform.InstanceDiff{
		Attributes: map[string]*terraform.ResourceAttrDiff{
			tc.changedKey: {Old: "before", New: "after"},
		},
		RawConfig: retriggerRawConfig(tc.configSet, tc.configVal, tc.changedKey),
	}
	if tc.configSet {
		diff.Attributes["retrigger_timeout_minutes"] = &terraform.ResourceAttrDiff{
			Old: tc.stateValue,
			New: fmt.Sprintf("%d", tc.configVal),
		}
	}
	data, err := schema.InternalMap(s.Schema).Data(&terraform.InstanceState{ID: "1", Attributes: attrs}, diff)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func retriggerCreateData(t *testing.T, s *schema.Resource, attrs map[string]*terraform.ResourceAttrDiff, set bool, val int64) *schema.ResourceData {
	t.Helper()
	rawAttrs := map[string]cty.Value{
		"retrigger_timeout_minutes": cty.NullVal(cty.Number),
	}
	if set {
		rawAttrs["retrigger_timeout_minutes"] = cty.NumberIntVal(val)
		attrs["retrigger_timeout_minutes"] = &terraform.ResourceAttrDiff{New: fmt.Sprintf("%d", val)}
	}
	diff := &terraform.InstanceDiff{
		Attributes: attrs,
		RawConfig:  cty.ObjectVal(rawAttrs),
	}
	data, err := schema.InternalMap(s.Schema).Data(nil, diff)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// captureRetriggerBody serves the JSON:API endpoints a resource hits during
// create/update+read and records the attributes hash of the write request.
func captureRetriggerBody(t *testing.T, createPath string, memberPath string, objType string) (*httptest.Server, *map[string]interface{}) {
	t.Helper()
	var got map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.api+json")
		if r.Method == http.MethodPost && r.URL.Path == createPath ||
			r.Method == http.MethodPut && r.URL.Path == memberPath {
			var body struct {
				Data struct {
					Attributes map[string]interface{} `json:"attributes"`
				} `json:"data"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Errorf("decoding %s request body: %v", r.Method, err)
				w.WriteHeader(400)
				return
			}
			got = body.Data.Attributes
		} else if r.Method != http.MethodGet {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"id": "1", "type": objType, "attributes": map[string]interface{}{},
			},
		})
	}))
	t.Cleanup(server.Close)
	return server, &got
}

func assertRetriggerAttr(t *testing.T, attrs map[string]interface{}, want interface{}) {
	t.Helper()
	got, present := attrs["retrigger_timeout_minutes"]
	if !present {
		t.Fatalf("request attributes missing retrigger_timeout_minutes: %#v", attrs)
	}
	if want == nil {
		if got != nil {
			t.Fatalf("expected retrigger_timeout_minutes null, got %#v", got)
		}
		return
	}
	if got != float64(want.(int64)) {
		t.Fatalf("expected retrigger_timeout_minutes %v, got %#v", want, got)
	}
}

func TestAlertUrgencyRetriggerTimeoutMinutesOnWire(t *testing.T) {
	resource := resourceAlertUrgency()
	server, got := captureRetriggerBody(t, "/v1/alert_urgencies", "/v1/alert_urgencies/1", "alert_urgencies")
	api, err := client.NewClient(server.URL, "test-token", "retrigger-test")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	for _, tc := range retriggerCases() {
		t.Run(tc.name, func(t *testing.T) {
			*got = nil
			data := retriggerUpdateData(t, resource, map[string]string{"name": "before"}, tc)
			if diags := resourceAlertUrgencyUpdate(ctx, data, api); diags.HasError() {
				t.Fatal(diags)
			}
			assertRetriggerAttr(t, *got, tc.want)
		})
	}

	t.Run("create config absent", func(t *testing.T) {
		*got = nil
		data := retriggerCreateData(t, resource, map[string]*terraform.ResourceAttrDiff{
			"name":    {New: "before"},
			"urgency": {New: "high"},
		}, false, 0)
		if diags := resourceAlertUrgencyCreate(ctx, data, api); diags.HasError() {
			t.Fatal(diags)
		}
		assertRetriggerAttr(t, *got, nil)
	})

	t.Run("create config -1", func(t *testing.T) {
		*got = nil
		data := retriggerCreateData(t, resource, map[string]*terraform.ResourceAttrDiff{
			"name":    {New: "before"},
			"urgency": {New: "high"},
		}, true, -1)
		if diags := resourceAlertUrgencyCreate(ctx, data, api); diags.HasError() {
			t.Fatal(diags)
		}
		assertRetriggerAttr(t, *got, int64(-1))
	})
}

func TestEscalationPathRetriggerTimeoutMinutesOnWire(t *testing.T) {
	resource := resourceEscalationPath()
	server, got := captureRetriggerBody(t, "/v1/escalation_policies/7/escalation_paths", "/v1/escalation_paths/1", "escalation_paths")
	api, err := client.NewClient(server.URL, "test-token", "retrigger-test")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	extraState := map[string]string{"name": "before"}

	for _, tc := range retriggerCases() {
		t.Run(tc.name, func(t *testing.T) {
			*got = nil
			data := retriggerUpdateData(t, resource, extraState, tc)
			if diags := resourceEscalationPathUpdate(ctx, data, api); diags.HasError() {
				t.Fatal(diags)
			}
			assertRetriggerAttr(t, *got, tc.want)
		})
	}

	t.Run("create config absent", func(t *testing.T) {
		*got = nil
		data := retriggerCreateData(t, resource, map[string]*terraform.ResourceAttrDiff{
			"name":                 {New: "before"},
			"escalation_policy_id": {New: "7"},
		}, false, 0)
		if diags := resourceEscalationPathCreate(ctx, data, api); diags.HasError() {
			t.Fatal(diags)
		}
		assertRetriggerAttr(t, *got, nil)
	})

	t.Run("create config -1", func(t *testing.T) {
		*got = nil
		data := retriggerCreateData(t, resource, map[string]*terraform.ResourceAttrDiff{
			"name":                 {New: "before"},
			"escalation_policy_id": {New: "7"},
		}, true, -1)
		if diags := resourceEscalationPathCreate(ctx, data, api); diags.HasError() {
			t.Fatal(diags)
		}
		assertRetriggerAttr(t, *got, int64(-1))
	})
}
