package stateupgrade

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// workflowTaskV0 mirrors the common workflow-task schema used through
// provider v5.17.x. Task-specific parameter schemas passed to it must remain
// frozen so Terraform can decode state written by those releases.
func workflowTaskV0(taskType string, taskParams map[string]*schema.Schema) *schema.Resource {
	taskParams["task_type"] = &schema.Schema{Type: schema.TypeString, Optional: true, Default: taskType}

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
				Elem:     &schema.Resource{Schema: taskParams},
			},
		},
	}
}

func optionalStringList() *schema.Schema {
	return &schema.Schema{Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}}
}

func namedObjectList() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{Schema: map[string]*schema.Schema{
			"id":   {Type: schema.TypeString, Required: true},
			"name": {Type: schema.TypeString, Required: true},
		}},
	}
}

// WorkflowTaskCreateGoogleCalendarEventV0 mirrors provider v5.17.x.
func WorkflowTaskCreateGoogleCalendarEventV0() *schema.Resource {
	return workflowTaskV0("create_google_calendar_event", map[string]*schema.Schema{
		"attendees":                   optionalStringList(),
		"time_zone":                   {Type: schema.TypeString, Optional: true},
		"calendar_id":                 {Type: schema.TypeString, Optional: true},
		"days_until_meeting":          {Type: schema.TypeString, Required: true},
		"time_of_meeting":             {Type: schema.TypeString, Required: true},
		"meeting_duration":            {Type: schema.TypeString, Required: true},
		"send_updates":                {Type: schema.TypeBool, Optional: true},
		"can_guests_modify_event":     {Type: schema.TypeBool, Optional: true},
		"can_guests_see_other_guests": {Type: schema.TypeBool, Optional: true},
		"can_guests_invite_others":    {Type: schema.TypeBool, Optional: true},
		"summary":                     {Type: schema.TypeString, Required: true},
		"description":                 {Type: schema.TypeString, Required: true},
		"exclude_weekends":            {Type: schema.TypeBool, Optional: true},
		"conference_solution_key":     {Type: schema.TypeString, Optional: true},
		"post_to_incident_timeline":   {Type: schema.TypeBool, Optional: true},
		"post_to_slack_channels":      namedObjectList(),
	})
}

// WorkflowTaskCreateMistralChatCompletionV0 mirrors provider v5.17.x.
func WorkflowTaskCreateMistralChatCompletionV0() *schema.Resource {
	return workflowTaskV0("create_mistral_chat_completion", map[string]*schema.Schema{
		"model":         {Type: schema.TypeMap, Required: true},
		"system_prompt": {Type: schema.TypeString, Optional: true},
		"prompt":        {Type: schema.TypeString, Required: true},
		"temperature":   {Type: schema.TypeInt, Optional: true},
		"max_tokens":    {Type: schema.TypeString, Optional: true},
		"top_p":         {Type: schema.TypeInt, Optional: true},
	})
}

// WorkflowTaskCreateOpenaiChatCompletionV0 mirrors provider v5.17.x.
func WorkflowTaskCreateOpenaiChatCompletionV0() *schema.Resource {
	return workflowTaskV0("create_openai_chat_completion", map[string]*schema.Schema{
		"model":             {Type: schema.TypeMap, Required: true},
		"system_prompt":     {Type: schema.TypeString, Optional: true},
		"prompt":            {Type: schema.TypeString, Required: true},
		"temperature":       {Type: schema.TypeInt, Optional: true},
		"max_tokens":        {Type: schema.TypeString, Optional: true},
		"top_p":             {Type: schema.TypeInt, Optional: true},
		"reasoning_effort":  {Type: schema.TypeString, Optional: true},
		"reasoning_summary": {Type: schema.TypeString, Optional: true},
	})
}

// WorkflowTaskCreateOutlookEventV0 mirrors provider v5.17.x.
func WorkflowTaskCreateOutlookEventV0() *schema.Resource {
	return workflowTaskV0("create_outlook_event", map[string]*schema.Schema{
		"calendar":                  {Type: schema.TypeMap, Required: true},
		"attendees":                 optionalStringList(),
		"time_zone":                 {Type: schema.TypeString, Optional: true},
		"days_until_meeting":        {Type: schema.TypeString, Required: true},
		"time_of_meeting":           {Type: schema.TypeString, Required: true},
		"meeting_duration":          {Type: schema.TypeString, Required: true},
		"summary":                   {Type: schema.TypeString, Required: true},
		"description":               {Type: schema.TypeString, Required: true},
		"exclude_weekends":          {Type: schema.TypeBool, Optional: true},
		"enable_online_meeting":     {Type: schema.TypeBool, Optional: true},
		"post_to_incident_timeline": {Type: schema.TypeBool, Optional: true},
		"post_to_slack_channels":    namedObjectList(),
	})
}

// WorkflowTaskUpdateGoogleCalendarEventV0 mirrors provider v5.17.x.
func WorkflowTaskUpdateGoogleCalendarEventV0() *schema.Resource {
	return workflowTaskV0("update_google_calendar_event", map[string]*schema.Schema{
		"calendar_id":                 {Type: schema.TypeString, Optional: true},
		"event_id":                    {Type: schema.TypeString, Required: true},
		"summary":                     {Type: schema.TypeString, Optional: true},
		"description":                 {Type: schema.TypeString, Optional: true},
		"adjustment_days":             {Type: schema.TypeString, Optional: true},
		"time_of_meeting":             {Type: schema.TypeString, Optional: true},
		"meeting_duration":            {Type: schema.TypeString, Optional: true},
		"send_updates":                {Type: schema.TypeBool, Optional: true},
		"can_guests_modify_event":     {Type: schema.TypeBool, Optional: true},
		"can_guests_see_other_guests": {Type: schema.TypeBool, Optional: true},
		"can_guests_invite_others":    {Type: schema.TypeBool, Optional: true},
		"attendees":                   optionalStringList(),
		"replace_attendees":           {Type: schema.TypeBool, Optional: true},
		"conference_solution_key":     {Type: schema.TypeString, Optional: true},
		"post_to_incident_timeline":   {Type: schema.TypeBool, Optional: true},
		"post_to_slack_channels":      namedObjectList(),
	})
}

// WorkflowTaskUpdatePagerdutyIncidentV0 mirrors provider v5.17.x.
func WorkflowTaskUpdatePagerdutyIncidentV0() *schema.Resource {
	return workflowTaskV0("update_pagerduty_incident", map[string]*schema.Schema{
		"pagerduty_incident_id": {Type: schema.TypeString, Required: true},
		"title":                 {Type: schema.TypeString, Optional: true},
		"status":                {Type: schema.TypeString, Optional: true},
		"resolution":            {Type: schema.TypeString, Optional: true},
		"escalation_level":      {Type: schema.TypeString, Optional: true},
		"urgency":               {Type: schema.TypeString, Optional: true},
		"priority":              {Type: schema.TypeString, Optional: true},
	})
}
