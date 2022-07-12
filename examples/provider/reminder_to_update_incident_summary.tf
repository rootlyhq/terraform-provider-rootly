resource "rootly_workflow_incident" "reminder_to_update_incident_summary" {
  name = "Reminder to update incident summary"
  description = "Reminds users after incident creation to update the incident summary."
  trigger_params {
    triggers = ["incident_created"]
    wait = "2 mins"
    incident_condition_summary = "UNSET"
  }
  enabled = true
}

resource "rootly_workflow_task_send_slack_message" "send_slack_message" {
  workflow_id = rootly_workflow_incident.reminder_to_update_incident_summary.id
  task_params {
    name = "Summary update reminder"
    channels {
      id = "{{ incident.slack_channel_id }}"
      name = "{{ incident.slack_channel_id }}"
    }
    actionables = ["update_summary"]
    text = "The incident is missing a summary."
  }
}