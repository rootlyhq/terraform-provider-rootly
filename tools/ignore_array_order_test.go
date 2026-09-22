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
