resource "rootly_workflow_incident" "auto_archive_incident" {
  name        = "Auto archive incident 48hrs after resolution"
  description = "After an incident has been resolved, automatically archive the channel after 48 hours."
  trigger_params {
    triggers                  = ["status_updated"]
    wait                      = "2 days"
    incident_statuses         = ["resolved"]
    incident_condition_status = "IS"
  }
  enabled = true
}

resource "rootly_workflow_task_archive_slack_channels" "archive_slack_channels" {
  workflow_id = rootly_workflow_incident.auto_archive_incident.id
  task_params {
    name = "Archive Slack channels"
    channels {
      id   = "{{ incident.slack_channel_id }}"
      name = "{{ incident.slack_channel_id }}"
    }
  }
}
