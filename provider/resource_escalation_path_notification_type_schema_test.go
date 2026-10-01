package provider

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v5/client"
	"github.com/stretchr/testify/require"
)

func TestResourceEscalationPathNotificationTypeSchema(t *testing.T) {
	t.Parallel()

	resourceSchema := resourceEscalationPath().Schema

	rules, ok := resourceSchema["notification_type_rules"]
	require.True(t, ok, "notification_type_rules should be part of the resource schema")
	require.Equal(t, schema.TypeList, rules.Type)
	require.True(t, rules.Optional)
	require.Equal(t, 10, rules.MaxItems)
	require.Nil(t, rules.DiffSuppressFunc, "rule order is semantic: first match wins")

	ruleSchema := rules.Elem.(*schema.Resource).Schema

	conditions, ok := ruleSchema["conditions"]
	require.True(t, ok, "conditions should be part of the notification_type_rules schema")
	require.True(t, conditions.Required)
	require.Equal(t, 1, conditions.MinItems)
	require.Equal(t, 5, conditions.MaxItems)
	require.NotNil(t, conditions.DiffSuppressFunc, "conditions are combined per match_mode")

	conditionSchema := conditions.Elem.(*schema.Resource).Schema
	for name := range escalationPathConditionSchema() {
		require.Contains(t, conditionSchema, name, "conditions should expose the same fields as path rules")
	}

	fallback, ok := resourceSchema["notification_type_fallback"]
	require.True(t, ok, "notification_type_fallback should be part of the resource schema")
	require.Equal(t, schema.TypeString, fallback.Type)
	require.True(t, fallback.Optional)

	for _, valid := range []string{"audible", "quiet"} {
		_, errs := fallback.ValidateFunc(valid, "notification_type_fallback")
		require.Empty(t, errs, "%q should be accepted", valid)

		_, errs = ruleSchema["notification_type"].ValidateFunc(valid, "notification_type")
		require.Empty(t, errs, "%q should be accepted", valid)
	}

	_, errs := fallback.ValidateFunc("loud", "notification_type_fallback")
	require.NotEmpty(t, errs, "unknown notification types should be rejected")

	_, errs = ruleSchema["notification_type"].ValidateFunc("loud", "notification_type")
	require.NotEmpty(t, errs, "unknown notification types should be rejected")

	_, errs = ruleSchema["match_mode"].ValidateFunc("match-some-rules", "match_mode")
	require.NotEmpty(t, errs, "unknown match modes should be rejected")
}

func TestResourceEscalationPathNotificationTypeRulesPayload(t *testing.T) {
	t.Parallel()

	resourceData := schema.TestResourceDataRaw(t, resourceEscalationPath().Schema, nil)
	require.NoError(t, resourceData.Set("notification_type_fallback", "quiet"))
	require.NoError(t, resourceData.Set("notification_type_rules", notificationTypeRulesConfig()))

	s := &client.EscalationPath{}
	if value, ok := resourceData.GetOkExists("notification_type_rules"); ok {
		s.NotificationTypeRules = value.([]interface{})
	}
	if value, ok := resourceData.GetOkExists("notification_type_fallback"); ok {
		s.NotificationTypeFallback = value.(string)
	}

	require.Len(t, s.NotificationTypeRules, 2)

	buffer, err := client.MarshalData(s)
	require.NoError(t, err)
	body, err := io.ReadAll(buffer)
	require.NoError(t, err)

	var payload struct {
		Data struct {
			Attributes struct {
				NotificationTypeFallback string `json:"notification_type_fallback"`
				NotificationTypeRules    []struct {
					NotificationType string                   `json:"notification_type"`
					MatchMode        string                   `json:"match_mode"`
					Conditions       []map[string]interface{} `json:"conditions"`
				} `json:"notification_type_rules"`
			} `json:"attributes"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(body, &payload))

	attributes := payload.Data.Attributes
	require.Equal(t, "quiet", attributes.NotificationTypeFallback)
	require.Len(t, attributes.NotificationTypeRules, 2)

	first := attributes.NotificationTypeRules[0]
	require.Equal(t, "audible", first.NotificationType)
	require.Equal(t, "match-all-rules", first.MatchMode)
	require.Len(t, first.Conditions, 2)
	require.Equal(t, "alert_urgency", first.Conditions[0]["rule_type"])
	require.Equal(t, []interface{}{"urgency-1"}, first.Conditions[0]["urgency_ids"])
	require.Equal(t, "deferral_window", first.Conditions[1]["rule_type"])
	require.Equal(t, "America/New_York", first.Conditions[1]["time_zone"])

	timeBlocks, ok := first.Conditions[1]["time_blocks"].([]interface{})
	require.True(t, ok, "time_blocks should serialize as a list")
	require.Len(t, timeBlocks, 1)
	timeBlock := timeBlocks[0].(map[string]interface{})
	require.Equal(t, "09:00", timeBlock["start_time"])
	require.Equal(t, "17:00", timeBlock["end_time"])
	require.Equal(t, true, timeBlock["monday"])
	require.Equal(t, false, timeBlock["sunday"])

	second := attributes.NotificationTypeRules[1]
	require.Equal(t, "quiet", second.NotificationType)
	require.Equal(t, "match-any-rule", second.MatchMode)
	require.Len(t, second.Conditions, 1)
	require.Equal(t, "json_path", second.Conditions[0]["rule_type"])
	require.Equal(t, "$.severity", second.Conditions[0]["json_path"])
	require.Equal(t, "info", second.Conditions[0]["value"])
}

func TestProcessEscalationPathNotificationTypeRulesPreservesOrder(t *testing.T) {
	t.Parallel()

	processed := processEscalationPathNotificationTypeRules(notificationTypeRulesAPIResponse())
	require.Len(t, processed, 3)

	require.Equal(t, "audible", processed[0]["notification_type"])
	require.Equal(t, "match-all-rules", processed[0]["match_mode"])
	require.Equal(t, "quiet", processed[1]["notification_type"])
	require.Equal(t, "audible", processed[2]["notification_type"])
	require.Equal(t, "match-any-rule", processed[2]["match_mode"])

	firstConditions := processed[0]["conditions"].([]map[string]interface{})
	require.Len(t, firstConditions, 2)
	require.Equal(t, "alert_urgency", firstConditions[0]["rule_type"])
	require.Equal(t, []interface{}{"urgency-1"}, firstConditions[0]["urgency_ids"])
	require.NotContains(t, firstConditions[0], "service_ids", "empty collections should stay unset")
	require.Equal(t, "deferral_window", firstConditions[1]["rule_type"])
	require.Equal(t, "America/New_York", firstConditions[1]["time_zone"])

	timeBlocks := firstConditions[1]["time_blocks"].([]map[string]interface{})
	require.Len(t, timeBlocks, 2)
	require.Equal(t, "09:00", timeBlocks[0]["start_time"])
	require.Equal(t, true, timeBlocks[1]["all_day"])

	secondConditions := processed[1]["conditions"].([]map[string]interface{})
	require.Len(t, secondConditions, 1)
	require.Equal(t, "json_path", secondConditions[0]["rule_type"])
	require.Equal(t, "info", secondConditions[0]["value"])

	require.Empty(t, processEscalationPathNotificationTypeRules([]interface{}{}))
}

func TestResourceEscalationPathNotificationTypeRulesRoundTrip(t *testing.T) {
	t.Parallel()

	resourceData := schema.TestResourceDataRaw(t, resourceEscalationPath().Schema, nil)
	require.NoError(t, resourceData.Set("notification_type_rules", processEscalationPathNotificationTypeRules(notificationTypeRulesAPIResponse())))

	stored := resourceData.Get("notification_type_rules").([]interface{})
	require.Len(t, stored, 3)

	types := make([]string, 0, len(stored))
	for _, rule := range stored {
		types = append(types, rule.(map[string]interface{})["notification_type"].(string))
	}
	require.Equal(t, []string{"audible", "quiet", "audible"}, types, "rule order drives which rule wins")

	first := stored[0].(map[string]interface{})
	require.Equal(t, "match-all-rules", first["match_mode"])

	conditions := first["conditions"].([]interface{})
	require.Len(t, conditions, 2)
	require.Equal(t, "alert_urgency", conditions[0].(map[string]interface{})["rule_type"])
	require.Equal(t, "deferral_window", conditions[1].(map[string]interface{})["rule_type"])
	require.Equal(t, "America/New_York", conditions[1].(map[string]interface{})["time_zone"])
}

func notificationTypeRulesConfig() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"notification_type": "audible",
			"match_mode":        "match-all-rules",
			"conditions": []interface{}{
				map[string]interface{}{
					"rule_type":   "alert_urgency",
					"urgency_ids": []interface{}{"urgency-1"},
				},
				map[string]interface{}{
					"rule_type": "deferral_window",
					"time_zone": "America/New_York",
					"time_blocks": []interface{}{
						map[string]interface{}{
							"monday":     true,
							"tuesday":    true,
							"wednesday":  true,
							"thursday":   true,
							"friday":     true,
							"start_time": "09:00",
							"end_time":   "17:00",
						},
					},
				},
			},
		},
		map[string]interface{}{
			"notification_type": "quiet",
			"match_mode":        "match-any-rule",
			"conditions": []interface{}{
				map[string]interface{}{
					"rule_type": "json_path",
					"json_path": "$.severity",
					"operator":  "is",
					"value":     "info",
				},
			},
		},
	}
}

func notificationTypeRulesAPIResponse() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"notification_type": "audible",
			"match_mode":        "match-all-rules",
			"conditions": []interface{}{
				map[string]interface{}{
					"rule_type":   "alert_urgency",
					"urgency_ids": []interface{}{"urgency-1"},
					"values":      []interface{}{},
					"service_ids": []interface{}{},
				},
				map[string]interface{}{
					"rule_type": "deferral_window",
					"time_zone": "America/New_York",
					"time_blocks": []interface{}{
						map[string]interface{}{
							"monday":     true,
							"start_time": "09:00",
							"end_time":   "17:00",
							"all_day":    false,
							"position":   1,
						},
						map[string]interface{}{
							"saturday": true,
							"sunday":   true,
							"all_day":  true,
							"position": 2,
						},
					},
				},
			},
		},
		map[string]interface{}{
			"notification_type": "quiet",
			"match_mode":        "match-all-rules",
			"conditions": []interface{}{
				map[string]interface{}{
					"rule_type": "json_path",
					"json_path": "$.severity",
					"operator":  "is",
					"value":     "info",
				},
			},
		},
		map[string]interface{}{
			"notification_type": "audible",
			"match_mode":        "match-any-rule",
			"conditions": []interface{}{
				map[string]interface{}{
					"rule_type": "related_incidents",
					"operator":  "is_set",
				},
			},
		},
	}
}
