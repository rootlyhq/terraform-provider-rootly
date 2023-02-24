resource "rootly_workflow_incident" "reminder_to_leave_incident_feedback" {
  name        = "Reminder to leave incident feedback"
  description = "Reminds users to leave feedback on an incident when it has been resolved or mitigated"
  trigger_params {
    triggers                  = ["status_updated"]
    incident_condition_status = "IS"
    incident_statuses         = ["resolved"]
  }
  enabled = true
}

resource "rootly_workflow_task_send_slack_message" "send_slack_message" {
  workflow_id = rootly_workflow_incident.reminder_to_leave_incident_feedback.id
  task_params {
    name = "Reminds users to leave incident feedback"
    channels {
      id   = "{{ incident.slack_channel_id }}"
      name = "{{ incident.slack_channel_id }}"
    }
    actionables = ["leave_feedback"]
    text        = "How did the incident go?"
  }
}
