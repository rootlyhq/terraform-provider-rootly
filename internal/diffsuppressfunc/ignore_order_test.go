package diffsuppressfunc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestEqualIgnoringOrderOmittingFields(t *testing.T) {
	omitNameWhenFallback := func(item map[string]interface{}) []string {
		if fallback, ok := item["fallback"].(bool); ok && fallback {
			return []string{"name"}
		}
		return nil
	}
	testSchema := map[string]*schema.Schema{
		"items": {
			Type:             schema.TypeList,
			Optional:         true,
			DiffSuppressFunc: EqualIgnoringOrderOmittingFields(omitNameWhenFallback),
			Elem: &schema.Resource{Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"fallback": {
					Type:     schema.TypeBool,
					Optional: true,
				},
			}},
		},
	}
	oldItems := []interface{}{
		map[string]interface{}{"name": "server name", "fallback": true},
		map[string]interface{}{"name": "rule", "fallback": false},
	}

	tests := []struct {
		name     string
		newItems []interface{}
		wantDiff bool
	}{
		{
			name: "omitted field differs",
			newItems: []interface{}{
				map[string]interface{}{"name": "configured name", "fallback": true},
				map[string]interface{}{"name": "rule", "fallback": false},
			},
			wantDiff: false,
		},
		{
			name: "omitted field differs and reordered",
			newItems: []interface{}{
				map[string]interface{}{"name": "rule", "fallback": false},
				map[string]interface{}{"name": "configured name", "fallback": true},
			},
			wantDiff: false,
		},
		{
			name: "field kept for other items",
			newItems: []interface{}{
				map[string]interface{}{"name": "server name", "fallback": true},
				map[string]interface{}{"name": "renamed rule", "fallback": false},
			},
			wantDiff: true,
		},
		{
			name: "item added",
			newItems: []interface{}{
				map[string]interface{}{"name": "server name", "fallback": true},
				map[string]interface{}{"name": "rule", "fallback": false},
				map[string]interface{}{"name": "another rule", "fallback": false},
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
