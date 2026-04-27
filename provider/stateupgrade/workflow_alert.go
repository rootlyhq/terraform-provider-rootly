// Package stateupgrade holds Terraform SDK StateUpgraders for hand-maintained
// resources. Each file here owns one resource's frozen prior-version schema(s)
// plus the upgrade function(s) that migrate stored state forward.
package stateupgrade

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// WorkflowAlertV0 mirrors the provider-v5.5.x schema shape for
// rootly_workflow_alert. It exists only so the SDK can decode pre-v5.6.0 state
// for UpgradeWorkflowAlertV0ToV1. Do not update this — it is frozen at 5.5.x.
func WorkflowAlertV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name":                     {Type: schema.TypeString, Required: true},
			"slug":                     {Type: schema.TypeString, Optional: true, Computed: true},
			"description":              {Type: schema.TypeString, Optional: true, Computed: true},
			"command":                  {Type: schema.TypeString, Optional: true, Computed: true},
			"command_feedback_enabled": {Type: schema.TypeBool, Optional: true, Computed: true},
			"wait":                     {Type: schema.TypeString, Optional: true, Computed: true},
			"repeat_every_duration":    {Type: schema.TypeString, Optional: true, Computed: true},
			"repeat_condition_duration_since_first_run": {Type: schema.TypeString, Optional: true, Computed: true},
			"repeat_condition_number_of_repeats":        {Type: schema.TypeInt, Optional: true, Computed: true},
			"continuously_repeat":                       {Type: schema.TypeBool, Optional: true, Computed: true},
			"repeat_on": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enabled":           {Type: schema.TypeBool, Optional: true},
			"locked":            {Type: schema.TypeBool, Optional: true, Computed: true},
			"position":          {Type: schema.TypeInt, Optional: true, Computed: true},
			"workflow_group_id": {Type: schema.TypeString, Optional: true, Computed: true},
			"trigger_params": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"trigger_type": {Type: schema.TypeString, Optional: true},
						"triggers": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"alert_condition":                    {Type: schema.TypeString, Optional: true},
						"alert_condition_source":             {Type: schema.TypeString, Optional: true},
						"alert_condition_source_use_regexp":  {Type: schema.TypeBool, Optional: true, Computed: true},
						"alert_sources":                      {Type: schema.TypeList, Optional: true, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"alert_condition_label":              {Type: schema.TypeString, Optional: true},
						"alert_condition_label_use_regexp":   {Type: schema.TypeBool, Optional: true, Computed: true},
						"alert_condition_status":             {Type: schema.TypeString, Optional: true},
						"alert_condition_status_use_regexp":  {Type: schema.TypeBool, Optional: true, Computed: true},
						"alert_statuses":                     {Type: schema.TypeList, Optional: true, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"alert_labels":                       {Type: schema.TypeList, Optional: true, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"alert_condition_payload":            {Type: schema.TypeString, Optional: true},
						"alert_condition_payload_use_regexp": {Type: schema.TypeBool, Optional: true, Computed: true},
						"alert_payload":                      {Type: schema.TypeList, Optional: true, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"alert_query_payload":                {Type: schema.TypeString, Optional: true, Computed: true},
						"alert_field_conditions": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id":   {Type: schema.TypeString, Required: true},
									"name": {Type: schema.TypeString, Required: true},
								},
							},
						},
						"alert_payload_conditions": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"environment_ids":   {Type: schema.TypeList, Optional: true, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
			"severity_ids":      {Type: schema.TypeList, Optional: true, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
			"incident_type_ids": {Type: schema.TypeList, Optional: true, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
			"incident_role_ids": {Type: schema.TypeList, Optional: true, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
			"service_ids":       {Type: schema.TypeList, Optional: true, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
			"functionality_ids": {Type: schema.TypeList, Optional: true, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
			"group_ids":         {Type: schema.TypeList, Optional: true, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
			"cause_ids":         {Type: schema.TypeList, Optional: true, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
			"sub_status_ids":    {Type: schema.TypeList, Optional: true, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
		},
	}
}

// UpgradeWorkflowAlertV0ToV1 migrates Terraform state from the v5.5.x schema
// to v5.6.0 by coercing alert_payload_conditions into the nested-block shape.
// See CoerceAlertPayloadConditions for the coercion itself.
func UpgradeWorkflowAlertV0ToV1(ctx context.Context, rawState map[string]any, _ any) (map[string]any, error) {
	tps, ok := rawState["trigger_params"].([]any)
	if !ok || len(tps) == 0 {
		return rawState, nil
	}
	if tp, ok := tps[0].(map[string]any); ok {
		CoerceAlertPayloadConditions(ctx, tp)
		tps[0] = tp
		rawState["trigger_params"] = tps
	}
	return rawState, nil
}

// CoerceAlertPayloadConditions normalizes the alert_payload_conditions field
// inside a workflow_alert trigger_params map to the v5.6.0 nested-block shape,
// in place. It handles:
//
//   - Missing/empty/nil: no-op (empty map is deleted so d.Set doesn't error).
//   - Already v5.6.0 shape (a list): no-op.
//   - v5.5.x shape { logic, conditions: "<json string>" }: parsed into
//     [ { logic, conditions: [ {query, operator, values, use_regexp} ] } ].
//   - Malformed conditions JSON: logs a warning and emits an empty conditions
//     list rather than blocking the plan (IR-4120: unblock customers stuck on
//     legacy-shaped API responses).
//
// Used by both the V0->V1 state upgrader and resourceWorkflowAlertRead, since
// existing API rows still return the legacy shape for rows written by v5.5.x.
func CoerceAlertPayloadConditions(ctx context.Context, triggerParams map[string]any) {
	apc, present := triggerParams["alert_payload_conditions"]
	if !present || apc == nil {
		return
	}

	if _, isList := apc.([]any); isList {
		return
	}

	apcMap, ok := apc.(map[string]any)
	if !ok {
		return
	}
	if len(apcMap) == 0 {
		delete(triggerParams, "alert_payload_conditions")
		return
	}

	logic, _ := apcMap["logic"].(string)
	if logic == "" {
		logic = "ALL"
	}

	var conditions []any
	switch raw := apcMap["conditions"].(type) {
	case string:
		if raw != "" {
			var parsed []map[string]any
			if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
				tflog.Warn(ctx, "failed to parse legacy alert_payload_conditions.conditions JSON; emitting empty conditions list", map[string]any{
					"error": err.Error(),
					"raw":   raw,
				})
			} else {
				conditions = buildConditions(parsed)
			}
		}
	case []any:
		parsed := make([]map[string]any, 0, len(raw))
		for _, item := range raw {
			if m, ok := item.(map[string]any); ok {
				parsed = append(parsed, m)
			}
		}
		conditions = buildConditions(parsed)
	}

	triggerParams["alert_payload_conditions"] = []any{
		map[string]any{
			"logic":      logic,
			"conditions": conditions,
		},
	}
}

func buildConditions(parsed []map[string]any) []any {
	out := make([]any, 0, len(parsed))
	for _, c := range parsed {
		upgraded := map[string]any{
			"query":    stringOrEmpty(c["query"]),
			"operator": stringOrEmpty(c["operator"]),
		}
		if v, ok := c["values"].([]any); ok {
			upgraded["values"] = v
		}
		if v, ok := c["use_regexp"].(bool); ok {
			upgraded["use_regexp"] = v
		}
		out = append(out, upgraded)
	}
	return out
}

func stringOrEmpty(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
