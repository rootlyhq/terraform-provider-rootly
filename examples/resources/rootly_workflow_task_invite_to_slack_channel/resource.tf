resource "rootly_workflow_incident" "invite_to_incident_channel" {
  name        = "Invite Slack users and groups to incident channel"
  description = "Invite Slack User Groups e.g. @oncall-infra to the incident channel."
  trigger_params {
    triggers                  = ["incident_created"]
    incident_statuses         = ["started"]
    incident_condition_status = "IS"
  }
  enabled = true
}

resource "rootly_workflow_task_invite_to_slack_channel" "invite_to_slack_channel" {
  workflow_id = rootly_workflow_incident.invite_to_incident_channel.id
  task_params {
    name = "Invite to Slack channel"
  }
  // TODO
}
