resource "rootly_workflow_incident" "send_message_in_incident_channel" {
  name        = "Send a message in the incident channel"
  description = "Send a message in the incident channel"
  trigger_params {
    triggers                  = ["incident_created"]
    incident_statuses         = ["started"]
    incident_condition_status = "IS"
  }
  enabled = true
}

resource "rootly_workflow_task_send_slack_blocks" "send_slack_blocks" {
  workflow_id     = rootly_workflow_incident.send_message_in_incident_channel.id
  skip_on_failure = false
  enabled         = true

  task_params {
    name    = "Send Slack message"
    message = ":boom: New incident!"
    channels {
      name = "{{ incident.slack_channel_id }}"
      id   = "{{ incident.slack_channel_id }}"
    }
    blocks = jsonencode(
      [
        {
          "text" = {
            "emoji" = true
            "text"  = "This is the incident title: {{ incident.title }}"
            "type"  = "plain_text"
          }
          "type" = "header"
        }
      ]
    )
  }
}
