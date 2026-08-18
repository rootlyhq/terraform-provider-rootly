package stateupgrade

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// WorkflowTaskHTTPClientV0 mirrors the schema used through provider v5.17.x.
// retry_count and retry_wait_time were strings in that schema; keep this
// frozen so Terraform can decode and upgrade state written by those releases.
func WorkflowTaskHTTPClientV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"workflow_id":     {Type: schema.TypeString, Required: true},
			"name":            {Type: schema.TypeString, Optional: true, Computed: true},
			"position":        {Type: schema.TypeInt, Optional: true, Computed: true},
			"skip_on_failure": {Type: schema.TypeBool, Optional: true},
			"enabled":         {Type: schema.TypeBool, Optional: true},
			"task_params": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_type":                 {Type: schema.TypeString, Optional: true},
						"headers":                   {Type: schema.TypeString, Optional: true},
						"params":                    {Type: schema.TypeString, Optional: true},
						"body":                      {Type: schema.TypeString, Optional: true},
						"url":                       {Type: schema.TypeString, Required: true},
						"event_url":                 {Type: schema.TypeString, Optional: true},
						"event_message":             {Type: schema.TypeString, Optional: true},
						"method":                    {Type: schema.TypeString, Optional: true},
						"succeed_on_status":         {Type: schema.TypeString, Required: true},
						"post_to_incident_timeline": {Type: schema.TypeBool, Optional: true},
						"post_to_slack_channels": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id":   {Type: schema.TypeString, Required: true},
									"name": {Type: schema.TypeString, Required: true},
								},
							},
						},
						"retry_count":     {Type: schema.TypeString, Optional: true},
						"retry_wait_time": {Type: schema.TypeString, Optional: true},
					},
				},
			},
		},
	}
}

// UpgradeWorkflowTaskHTTPClientV0ToV1 converts the two retry settings changed
// from strings to integers in provider v5.18.0.
func UpgradeWorkflowTaskHTTPClientV0ToV1(_ context.Context, rawState map[string]any, _ any) (map[string]any, error) {
	taskParams, ok := rawState["task_params"].([]any)
	if !ok || len(taskParams) == 0 {
		return rawState, nil
	}

	params, ok := taskParams[0].(map[string]any)
	if !ok {
		return rawState, nil
	}

	for _, field := range []string{"retry_count", "retry_wait_time"} {
		value, present := params[field]
		if !present || value == nil {
			continue
		}

		switch value := value.(type) {
		case int:
			continue
		case float64:
			continue
		case string:
			if value == "" {
				delete(params, field)
				continue
			}
			converted, err := strconv.Atoi(value)
			if err != nil {
				delete(params, field)
				continue
			}
			params[field] = converted
		default:
			delete(params, field)
		}
	}

	taskParams[0] = params
	rawState["task_params"] = taskParams
	return rawState, nil
}
