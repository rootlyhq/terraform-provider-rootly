resource "rootly_workflow_incident" "renaming_slack_channel_to_jira_ticket_number" {
  name        = "Rename Slack channel to Jira ticket slug"
  description = "Automatically renames Slack channel to attached Jira ticket number."
  trigger_params {
    triggers                  = ["incident_created"]
    incident_statuses         = ["started"]
    wait                      = "30 seconds"
    incident_condition_status = "IS"
  }
  enabled = true
}

resource "rootly_workflow_task_rename_slack_channel" "rename_slack_channel" {
  workflow_id = rootly_workflow_incident.renaming_slack_channel_to_jira_ticket_number.id
  task_params {
    name = "Rename a Slack channel to Jira ticket slug"
    channel = {
      id   = "{{ incident.slack_channel_id }}"
      name = "{{ incident.slack_channel_id }}"
    }
    title = "incident-{{ incident.jira_issue_key }}"
  }
}
