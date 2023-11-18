resource "rootly_workflow_incident" "reminder_to_assign_incident_roles" {
  name        = "Reminder to assign incident roles"
  description = "Reminds users to assign incident roles."
  trigger_params {
    triggers                          = ["incident_created"]
    wait                              = "2 mins"
    incident_condition_incident_roles = "UNSET"
  }
  enabled = true
}

resource "rootly_workflow_task_send_slack_message" "send_slack_message" {
  workflow_id = rootly_workflow_incident.reminder_to_assign_incident_roles.id
  name        = "Reminds users to assign incident roles"

  task_params {
    channels {
      id   = "{{ incident.slack_channel_id }}"
      name = "{{ incident.slack_channel_id }}"
    }
    actionables = ["manage_incident_roles"]
    text        = "Some incident roles are still unassigned."
  }
}
