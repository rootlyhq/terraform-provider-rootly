package stateupgrade

import (
	"context"
	"math"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// WorkflowTaskHTTPClientV0 mirrors the schema used through provider v5.17.x.
// retry_count and retry_wait_time were strings in that schema; keep this
// frozen so Terraform can decode and upgrade state written by those releases.
func WorkflowTaskHTTPClientV0() *schema.Resource {
	return workflowTaskV0("http_client", map[string]*schema.Schema{
		"headers":                   {Type: schema.TypeString, Optional: true},
		"params":                    {Type: schema.TypeString, Optional: true},
		"body":                      {Type: schema.TypeString, Optional: true},
		"url":                       {Type: schema.TypeString, Required: true},
		"event_url":                 {Type: schema.TypeString, Optional: true},
		"event_message":             {Type: schema.TypeString, Optional: true},
		"method":                    {Type: schema.TypeString, Optional: true},
		"succeed_on_status":         {Type: schema.TypeString, Required: true},
		"post_to_incident_timeline": {Type: schema.TypeBool, Optional: true},
		"post_to_slack_channels":    namedObjectList(),
		"retry_count":               {Type: schema.TypeString, Optional: true},
		"retry_wait_time":           {Type: schema.TypeString, Optional: true},
	})
}

// UpgradeWorkflowTaskHTTPClientV0ToV1 converts the two retry settings changed
// from strings to integers in provider v5.18.0.
func UpgradeWorkflowTaskHTTPClientV0ToV1(_ context.Context, rawState map[string]any, _ any) (map[string]any, error) {
	return upgradeWorkflowTaskFields(rawState, map[string]workflowTaskFieldType{
		"retry_count":     workflowTaskFieldInt,
		"retry_wait_time": workflowTaskFieldInt,
	}), nil
}

func UpgradeWorkflowTaskCreateGoogleCalendarEventV0ToV1(_ context.Context, rawState map[string]any, _ any) (map[string]any, error) {
	return upgradeWorkflowTaskFields(rawState, map[string]workflowTaskFieldType{"days_until_meeting": workflowTaskFieldInt}), nil
}

func UpgradeWorkflowTaskCreateMistralChatCompletionV0ToV1(_ context.Context, rawState map[string]any, _ any) (map[string]any, error) {
	return upgradeWorkflowTaskFields(rawState, map[string]workflowTaskFieldType{
		"temperature": workflowTaskFieldString,
		"max_tokens":  workflowTaskFieldInt,
		"top_p":       workflowTaskFieldString,
	}), nil
}

func UpgradeWorkflowTaskCreateOpenaiChatCompletionV0ToV1(_ context.Context, rawState map[string]any, _ any) (map[string]any, error) {
	return upgradeWorkflowTaskFields(rawState, map[string]workflowTaskFieldType{
		"temperature": workflowTaskFieldString,
		"max_tokens":  workflowTaskFieldInt,
		"top_p":       workflowTaskFieldString,
	}), nil
}

func UpgradeWorkflowTaskCreateOutlookEventV0ToV1(_ context.Context, rawState map[string]any, _ any) (map[string]any, error) {
	return upgradeWorkflowTaskFields(rawState, map[string]workflowTaskFieldType{"days_until_meeting": workflowTaskFieldInt}), nil
}

func UpgradeWorkflowTaskUpdateGoogleCalendarEventV0ToV1(_ context.Context, rawState map[string]any, _ any) (map[string]any, error) {
	return upgradeWorkflowTaskFields(rawState, map[string]workflowTaskFieldType{"adjustment_days": workflowTaskFieldInt}), nil
}

func UpgradeWorkflowTaskUpdatePagerdutyIncidentV0ToV1(_ context.Context, rawState map[string]any, _ any) (map[string]any, error) {
	return upgradeWorkflowTaskFields(rawState, map[string]workflowTaskFieldType{"escalation_level": workflowTaskFieldInt}), nil
}

type workflowTaskFieldType int

const (
	workflowTaskFieldInt workflowTaskFieldType = iota
	workflowTaskFieldString
)

func upgradeWorkflowTaskFields(rawState map[string]any, fields map[string]workflowTaskFieldType) map[string]any {
	taskParams, ok := rawState["task_params"].([]any)
	if !ok || len(taskParams) == 0 {
		return rawState
	}

	params, ok := taskParams[0].(map[string]any)
	if !ok {
		return rawState
	}

	for field, targetType := range fields {
		value, present := params[field]
		if !present || value == nil {
			continue
		}

		converted, ok := convertWorkflowTaskField(value, targetType)
		if !ok {
			delete(params, field)
			continue
		}
		params[field] = converted
	}

	taskParams[0] = params
	rawState["task_params"] = taskParams
	return rawState
}

func convertWorkflowTaskField(value any, targetType workflowTaskFieldType) (any, bool) {
	switch targetType {
	case workflowTaskFieldInt:
		switch value := value.(type) {
		case int:
			return value, true
		case float64:
			return value, math.Trunc(value) == value
		case string:
			if value == "" {
				return nil, false
			}
			converted, err := strconv.Atoi(value)
			return converted, err == nil
		}
	case workflowTaskFieldString:
		switch value := value.(type) {
		case string:
			return value, value != ""
		case int:
			return strconv.Itoa(value), true
		case float64:
			if math.Trunc(value) == value {
				return strconv.FormatInt(int64(value), 10), true
			}
		}
	}

	return nil, false
}
