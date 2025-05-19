resource "rootly_workflow_incident" "notify_slack_channels" {
  name        = "Notify teams on Slack about the incident"
  description = "Send a message to specific teams on Slack regarding the incident."
  trigger_params {
    triggers                  = ["incident_created"]
    incident_statuses         = ["started"]
    incident_condition_status = "IS"
  }
  enabled = true
}

# Import examples:
# 1. Using import block (Terraform v1.5.0 and later):
# import {
#   to rootly_workflow_task_send_slack_message.send_slack_message
#   id = "00000000-0000-0000-0000-000000000000"
# }
#
# 2. Using terraform import command:
# terraform import rootly_workflow_task_send_slack_message.send_slack_message 00000000-0000-0000-0000-000000000000

resource "rootly_workflow_task_send_slack_message" "send_slack_message" {
  workflow_id     = rootly_workflow_incident.notify_slack_channels.id
  skip_on_failure = false
  enabled         = true
  task_params {
    name = "Notify team about incident"
    channels {
      id   = "{{ incident.slack_channel_id }}"
      name = "{{ incident.slack_channel_id }}"
    }
    text = "Heads up - wanted to let your team know we have an active incident."
  }
}
