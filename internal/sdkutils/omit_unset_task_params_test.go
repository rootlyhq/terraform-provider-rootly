package sdkutils

import (
	"context"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func omitUnsetTaskParamsData(t *testing.T, sync cty.Value) *schema.ResourceData {
	t.Helper()
	resource := &schema.Resource{Schema: map[string]*schema.Schema{
		"task_params": {
			Type:     schema.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &schema.Resource{Schema: map[string]*schema.Schema{
				"sync":    {Type: schema.TypeBool, Optional: true, Computed: true},
				"enabled": {Type: schema.TypeBool, Optional: true},
			}},
		},
	}}
	params := map[string]interface{}{"enabled": true}
	if !sync.IsNull() {
		params["sync"] = sync.True()
	}
	config := terraform.NewResourceConfigRaw(map[string]interface{}{"task_params": []interface{}{params}})
	diff, err := resource.SimpleDiff(context.Background(), nil, config, nil)
	if err != nil {
		t.Fatalf("diff: %v", err)
	}
	diff.RawConfig = cty.ObjectVal(map[string]cty.Value{
		"task_params": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
			"sync":    sync,
			"enabled": cty.True,
		})}),
	})
	data, err := schema.InternalMap(resource.Schema).Data(nil, diff)
	if err != nil {
		t.Fatalf("data: %v", err)
	}
	return data
}

func TestOmitUnsetTaskParamsDropsParamsMissingFromConfig(t *testing.T) {
	data := omitUnsetTaskParamsData(t, cty.NullVal(cty.Bool))
	params := map[string]interface{}{"sync": false, "enabled": true}

	OmitUnsetTaskParams(data, params, "sync")

	if _, ok := params["sync"]; ok {
		t.Fatalf("expected unset sync to be omitted, got %v", params)
	}
	if params["enabled"] != true {
		t.Fatalf("expected other params to be kept, got %v", params)
	}
}

func TestOmitUnsetTaskParamsKeepsExplicitFalse(t *testing.T) {
	data := omitUnsetTaskParamsData(t, cty.False)
	params := map[string]interface{}{"sync": false, "enabled": true}

	OmitUnsetTaskParams(data, params, "sync")

	if value, ok := params["sync"]; !ok || value != false {
		t.Fatalf("expected explicit false to be sent, got %v", params)
	}
}
