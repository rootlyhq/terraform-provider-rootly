package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestResourceAlertRouteRulesDiff(t *testing.T) {
	alertRouteSchema := resourceAlertRoute().Schema

	route := func(rules ...interface{}) map[string]interface{} {
		return map[string]interface{}{
			"name":              "Example Route",
			"enabled":           true,
			"alerts_source_ids": []interface{}{"source-1"},
			"rules":             rules,
		}
	}
	rule := func(name string, position int, fallback bool, targetID string) map[string]interface{} {
		return map[string]interface{}{
			"name":          name,
			"position":      position,
			"fallback_rule": fallback,
			"destinations": []interface{}{
				map[string]interface{}{"target_type": "Group", "target_id": targetID},
			},
		}
	}

	tests := []struct {
		name     string
		state    map[string]interface{}
		config   map[string]interface{}
		wantDiff bool
	}{
		{
			name:     "fallback rule name assigned by Rootly",
			state:    route(rule("Fallback Rule for Example Route", 1, true, "team-1")),
			config:   route(rule("Page on-call", 1, true, "team-1")),
			wantDiff: false,
		},
		{
			name:     "fallback rule destination changed",
			state:    route(rule("Fallback Rule for Example Route", 1, true, "team-1")),
			config:   route(rule("Page on-call", 1, true, "team-2")),
			wantDiff: true,
		},
		{
			name:     "regular rule name changed",
			state:    route(rule("Old name", 1, false, "team-1")),
			config:   route(rule("New name", 1, false, "team-1")),
			wantDiff: true,
		},
		{
			name:     "rule changed from regular to fallback",
			state:    route(rule("Page on-call", 1, false, "team-1")),
			config:   route(rule("Page on-call", 1, true, "team-1")),
			wantDiff: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := schema.TestResourceDataRaw(t, alertRouteSchema, test.state)
			state.SetId("route-1")
			config := terraform.NewResourceConfigRaw(test.config)

			diff, err := schema.InternalMap(alertRouteSchema).Diff(context.Background(), state.State(), config, nil, nil, true)
			if err != nil {
				t.Fatal(err)
			}
			gotDiff := diff != nil && !diff.Empty()
			if gotDiff != test.wantDiff {
				t.Fatalf("diff=%t, want %t: %#v", gotDiff, test.wantDiff, diff)
			}
		})
	}
}
