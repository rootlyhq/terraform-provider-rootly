resource "rootly_workflow_incident" "reminder_to_update_status_page" {
  name = "Reminder to update status page"
  description = "Periodically reminds the incident channel to update status page every 30 min."
  trigger_params {
    triggers = ["incident_created"]
    repeat_every_duration = "30 mins"
    wait = "1 min"
    incident_statuses = ["started"]
    incident_condition_status = "IS"
  }
  enabled = true
}

resource "rootly_workflow_task_send_slack_message" "send_slack_message" {
  workflow_id = rootly_workflow_incident.reminder_to_update_status_page.id
  task_params {
    name = "Status Page Update Reminder"
    channels {
      id = "{{ incident.slack_channel_id }}"
      name = "{{ incident.slack_channel_id }}"
    }
    text = "Reminder to update status page. Your next reminder is in 30 mins."
  }
}