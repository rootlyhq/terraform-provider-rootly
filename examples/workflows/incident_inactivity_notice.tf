resource "rootly_workflow_incident" "incident_inactivity_notice" {
  name        = "Surface helper when inactive"
  description = "Surfaces helper actions if incident hasn't had any user activity for 5 mins"
  trigger_params {
    triggers                        = ["incident_updated"]
    wait                            = "5 mins"
    incident_conditional_inactivity = "IS"
    incident_inactivity_duration    = "5 mins"
  }
  enabled = true
}

resource "rootly_workflow_task_send_slack_message" "send_slack_message" {
  workflow_id = rootly_workflow_incident.incident_inactivity_notice.id
  task_params {
    name = "Helper actions for user"
    channels {
      id   = "{{ incident.slack_channel_id }}"
      name = "{{ incident.slack_channel_id }}"
    }
    actionables = ["update_summary", "manage_incident_roles", "update_incident", "all_commands"]
    text        = <<EOT
*Need help?*
We noticed there hasn't been any recent activity in the incident channel. Here are a few things you can try if you need help.
EOT
  }
}
