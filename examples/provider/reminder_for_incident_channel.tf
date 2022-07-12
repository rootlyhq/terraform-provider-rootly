resource "rootly_workflow_incident" "reminder_for_incident_channel" {
  name = "Reminder for Incident Channel"
  description = "Periodically remind the incident channel to do a certain task, e.g. update Summary."
  trigger_params {
    triggers = ["incident_created"]
    incident_statuses = ["started"]
    incident_condition_status = "IS"
  }
  enabled = true
}

resource "rootly_workflow_task_send_slack_message" "send_slack_message" {
  workflow_id = rootly_workflow_incident.reminder_for_incident_channel.id
  task_params {
    name = "Reminder to do X"
    channels {
      id = "{{ incident.slack_channel_id }}"
      name = "{{ incident.slack_channel_id }}"
    }
    text = "Reminder is your periodic reminder to do X."
  }
}