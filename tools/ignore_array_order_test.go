package tools

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestEqualIgnoringOrderResolvesNestedListPath(t *testing.T) {
	if got := listPathFromDiffKey("task_params.0.selected_component_keys.0"); got != "task_params.0.selected_component_keys" {
		t.Fatalf("unexpected nested list path: %q", got)
	}
	if !listsAreEqual(
		[]interface{}{"Service:1", "Service:2"},
		[]interface{}{"Service:2", "Service:1"},
	) {
		t.Fatal("expected an order-only nested list change to be suppressed")
	}
}

func TestEqualIgnoringOrderNestedObjectList(t *testing.T) {
	testSchema := map[string]*schema.Schema{
		"items": {
			Type:             schema.TypeList,
			Optional:         true,
			DiffSuppressFunc: EqualIgnoringOrder,
			Elem: &schema.Resource{Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Required: true,
				},
			}},
		},
	}
	oldItems := []interface{}{
		map[string]interface{}{"id": "a"},
		map[string]interface{}{"id": "b"},
		map[string]interface{}{"id": "c"},
	}

	tests := []struct {
		name     string
		newItems []interface{}
		wantDiff bool
	}{
		{
			name: "reordered",
			newItems: []interface{}{
				map[string]interface{}{"id": "b"},
				map[string]interface{}{"id": "c"},
				map[string]interface{}{"id": "a"},
			},
			wantDiff: false,
		},
		{
			name: "changed",
			newItems: []interface{}{
				map[string]interface{}{"id": "a"},
				map[string]interface{}{"id": "b"},
				map[string]interface{}{"id": "d"},
			},
			wantDiff: true,
		},
		{
			name: "added",
			newItems: []interface{}{
				map[string]interface{}{"id": "a"},
				map[string]interface{}{"id": "b"},
				map[string]interface{}{"id": "c"},
				map[string]interface{}{"id": "d"},
			},
			wantDiff: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			oldData := schema.TestResourceDataRaw(t, testSchema, map[string]interface{}{"items": oldItems})
			oldData.SetId("test")
			config := terraform.NewResourceConfigRaw(map[string]interface{}{"items": test.newItems})

			diff, err := schema.InternalMap(testSchema).Diff(context.Background(), oldData.State(), config, nil, nil, true)
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

func TestEqualIgnoringOrderAndFieldsIgnoresComputedFields(t *testing.T) {
	testSchema := map[string]*schema.Schema{
		"items": {
			Type:             schema.TypeList,
			Optional:         true,
			DiffSuppressFunc: EqualIgnoringOrderAndFields([]string{"name", "value"}),
			Elem: &schema.Resource{Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"value": {
					Type:     schema.TypeString,
					Required: true,
				},
			}},
		},
	}

	oldItems := []interface{}{
		map[string]interface{}{"id": "old-id-1", "name": "field1", "value": "val1"},
		map[string]interface{}{"id": "old-id-2", "name": "field2", "value": "val2"},
	}

	tests := []struct {
		name     string
		newItems []interface{}
		wantDiff bool
	}{
		{
			name: "same fields, different computed id (no diff)",
			newItems: []interface{}{
				map[string]interface{}{"id": "new-id-1", "name": "field1", "value": "val1"},
				map[string]interface{}{"id": "new-id-2", "name": "field2", "value": "val2"},
			},
			wantDiff: false,
		},
		{
			name: "reordered with different ids (no diff)",
			newItems: []interface{}{
				map[string]interface{}{"id": "different-id", "name": "field2", "value": "val2"},
				map[string]interface{}{"id": "another-id", "name": "field1", "value": "val1"},
			},
			wantDiff: false,
		},
		{
			name: "changed value (should diff)",
			newItems: []interface{}{
				map[string]interface{}{"id": "any-id", "name": "field1", "value": "val1"},
				map[string]interface{}{"id": "any-id", "name": "field2", "value": "different"},
			},
			wantDiff: true,
		},
		{
			name: "changed name (should diff)",
			newItems: []interface{}{
				map[string]interface{}{"id": "any-id", "name": "field1", "value": "val1"},
				map[string]interface{}{"id": "any-id", "name": "field3", "value": "val2"},
			},
			wantDiff: true,
		},
		{
			name: "missing id in new (no diff if name and value match)",
			newItems: []interface{}{
				map[string]interface{}{"name": "field1", "value": "val1"},
				map[string]interface{}{"name": "field2", "value": "val2"},
			},
			wantDiff: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			oldData := schema.TestResourceDataRaw(t, testSchema, map[string]interface{}{"items": oldItems})
			oldData.SetId("test")
			config := terraform.NewResourceConfigRaw(map[string]interface{}{"items": test.newItems})

			diff, err := schema.InternalMap(testSchema).Diff(context.Background(), oldData.State(), config, nil, nil, true)
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

func TestEqualIgnoringOrderOfListComparesNamedListForNestedKeys(t *testing.T) {
	testSchema := map[string]*schema.Schema{
		"rules": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{Schema: map[string]*schema.Schema{
				"destinations": {
					Type:             schema.TypeList,
					Optional:         true,
					DiffSuppressFunc: EqualIgnoringOrderOfList("destinations"),
					Elem: &schema.Resource{Schema: map[string]*schema.Schema{
						"target_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"match_keys": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{Schema: map[string]*schema.Schema{
								"source_ref": {
									Type:     schema.TypeString,
									Required: true,
								},
							}},
						},
					}},
				},
			}},
		},
	}
	static := map[string]interface{}{"target_id": "escalation-policy"}
	dynamic := func(sourceRef string) map[string]interface{} {
		return map[string]interface{}{"match_keys": []interface{}{map[string]interface{}{"source_ref": sourceRef}}}
	}
	rules := func(destinations ...interface{}) map[string]interface{} {
		return map[string]interface{}{"rules": []interface{}{map[string]interface{}{"destinations": destinations}}}
	}

	tests := []struct {
		name     string
		newRules map[string]interface{}
		wantDiff bool
	}{
		{
			name:     "dynamic destination listed before the static one",
			newRules: rules(dynamic("$.service"), static),
			wantDiff: false,
		},
		{
			name:     "changed match key",
			newRules: rules(static, dynamic("$.team")),
			wantDiff: true,
		},
		{
			name:     "removed destination",
			newRules: rules(dynamic("$.service")),
			wantDiff: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			oldData := schema.TestResourceDataRaw(t, testSchema, rules(static, dynamic("$.service")))
			oldData.SetId("test")
			config := terraform.NewResourceConfigRaw(test.newRules)

			diff, err := schema.InternalMap(testSchema).Diff(context.Background(), oldData.State(), config, nil, nil, true)
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
