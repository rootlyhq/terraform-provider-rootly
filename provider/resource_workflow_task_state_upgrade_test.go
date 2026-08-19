package provider

import (
	"context"
	"math/big"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-cty/cty/msgpack"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestWorkflowTaskV518StateUpgrades(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		resourceName string
		resource     *schema.Resource
		legacyState  string
		currentState string
		expected     map[string]any
	}{
		{
			name:         "HTTP client retries",
			resourceName: "rootly_workflow_task_http_client",
			resource:     resourceWorkflowTaskHttpClient(),
			legacyState:  `{"workflow_id":"workflow-id","task_params":[{"url":"https://example.com","succeed_on_status":"200","retry_count":"1","retry_wait_time":"15"}]}`,
			currentState: `{"workflow_id":"workflow-id","task_params":[{"url":"https://example.com","succeed_on_status":"200","retry_count":1,"retry_wait_time":15}]}`,
			expected:     map[string]any{"url": "https://example.com", "retry_count": int64(1), "retry_wait_time": int64(15)},
		},
		{
			name:         "Google Calendar create days",
			resourceName: "rootly_workflow_task_create_google_calendar_event",
			resource:     resourceWorkflowTaskCreateGoogleCalendarEvent(),
			legacyState:  `{"workflow_id":"workflow-id","task_params":[{"days_until_meeting":"2","time_of_meeting":"09:00","meeting_duration":"30 minutes","summary":"Summary","description":"Description"}]}`,
			currentState: `{"workflow_id":"workflow-id","task_params":[{"days_until_meeting":2,"time_of_meeting":"09:00","meeting_duration":"30 minutes","summary":"Summary","description":"Description"}]}`,
			expected:     map[string]any{"days_until_meeting": int64(2), "time_of_meeting": "09:00"},
		},
		{
			name:         "Mistral numeric settings",
			resourceName: "rootly_workflow_task_create_mistral_chat_completion",
			resource:     resourceWorkflowTaskCreateMistralChatCompletion(),
			legacyState:  `{"workflow_id":"workflow-id","task_params":[{"model":{"id":"model-id"},"prompt":"Prompt","temperature":1,"max_tokens":"512","top_p":1}]}`,
			currentState: `{"workflow_id":"workflow-id","task_params":[{"model":{"id":"model-id"},"prompt":"Prompt","temperature":"1","max_tokens":512,"top_p":"1"}]}`,
			expected:     map[string]any{"prompt": "Prompt", "temperature": "1", "max_tokens": int64(512), "top_p": "1"},
		},
		{
			name:         "OpenAI numeric settings",
			resourceName: "rootly_workflow_task_create_openai_chat_completion",
			resource:     resourceWorkflowTaskCreateOpenaiChatCompletion(),
			legacyState:  `{"workflow_id":"workflow-id","task_params":[{"model":{"id":"model-id"},"prompt":"Prompt","temperature":1,"max_tokens":"1024","top_p":1}]}`,
			currentState: `{"workflow_id":"workflow-id","task_params":[{"model":{"id":"model-id"},"prompt":"Prompt","temperature":"1","max_tokens":1024,"top_p":"1"}]}`,
			expected:     map[string]any{"prompt": "Prompt", "temperature": "1", "max_tokens": int64(1024), "top_p": "1"},
		},
		{
			name:         "Outlook create days",
			resourceName: "rootly_workflow_task_create_outlook_event",
			resource:     resourceWorkflowTaskCreateOutlookEvent(),
			legacyState:  `{"workflow_id":"workflow-id","task_params":[{"calendar":{"id":"calendar-id"},"days_until_meeting":"3","time_of_meeting":"09:00","meeting_duration":"30 minutes","summary":"Summary","description":"Description"}]}`,
			currentState: `{"workflow_id":"workflow-id","task_params":[{"calendar":{"id":"calendar-id"},"days_until_meeting":3,"time_of_meeting":"09:00","meeting_duration":"30 minutes","summary":"Summary","description":"Description"}]}`,
			expected:     map[string]any{"days_until_meeting": int64(3), "summary": "Summary"},
		},
		{
			name:         "Google Calendar update adjustment",
			resourceName: "rootly_workflow_task_update_google_calendar_event",
			resource:     resourceWorkflowTaskUpdateGoogleCalendarEvent(),
			legacyState:  `{"workflow_id":"workflow-id","task_params":[{"event_id":"event-id","adjustment_days":"4"}]}`,
			currentState: `{"workflow_id":"workflow-id","task_params":[{"event_id":"event-id","adjustment_days":4}]}`,
			expected:     map[string]any{"event_id": "event-id", "adjustment_days": int64(4)},
		},
		{
			name:         "PagerDuty escalation level",
			resourceName: "rootly_workflow_task_update_pagerduty_incident",
			resource:     resourceWorkflowTaskUpdatePagerdutyIncident(),
			legacyState:  `{"workflow_id":"workflow-id","task_params":[{"pagerduty_incident_id":"incident-id","escalation_level":"5"}]}`,
			currentState: `{"workflow_id":"workflow-id","task_params":[{"pagerduty_incident_id":"incident-id","escalation_level":5}]}`,
			expected:     map[string]any{"pagerduty_incident_id": "incident-id", "escalation_level": int64(5)},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			for stateName, rawState := range map[string]string{
				"pre-v5.18 state": testCase.legacyState,
				"v5.18 state":     testCase.currentState,
			} {
				rawState := rawState
				t.Run(stateName, func(t *testing.T) {
					t.Parallel()

					server := schema.NewGRPCProviderServer(&schema.Provider{
						ResourcesMap: map[string]*schema.Resource{testCase.resourceName: testCase.resource},
					})
					response, err := server.UpgradeResourceState(context.Background(), &tfprotov5.UpgradeResourceStateRequest{
						TypeName: testCase.resourceName,
						Version:  0,
						RawState: &tfprotov5.RawState{JSON: []byte(rawState)},
					})
					require.NoError(t, err)
					require.Empty(t, response.Diagnostics)

					state, err := msgpack.Unmarshal(response.UpgradedState.MsgPack, testCase.resource.CoreConfigSchema().ImpliedType())
					require.NoError(t, err)
					params := state.GetAttr("task_params").Index(cty.NumberIntVal(0))

					for field, expected := range testCase.expected {
						value := params.GetAttr(field)
						switch expected := expected.(type) {
						case int64:
							require.Equal(t, expected, ctyNumberAsInt64(t, value))
						case string:
							require.Equal(t, expected, value.AsString())
						default:
							t.Fatalf("unsupported expected value type %T", expected)
						}
					}
				})
			}
		})
	}
}

func ctyNumberAsInt64(t *testing.T, value cty.Value) int64 {
	t.Helper()

	result, accuracy := value.AsBigFloat().Int64()
	require.Equal(t, big.Exact, accuracy)
	return result
}
