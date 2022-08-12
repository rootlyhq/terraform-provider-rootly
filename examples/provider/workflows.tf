# Jira workflow
resource "rootly_workflow_incident" "jira" {
  name        = "Create a Jira Issue"
  description = "Open Jira ticket whenever incident starts"
  trigger_params {
    triggers                  = ["incident_created"]
    incident_condition_kind   = "IS"
    incident_kinds            = ["normal"]
    incident_condition_status = "IS"
    incident_statuses         = ["started"]
  }
  enabled = true
}

resource "rootly_workflow_task_create_jira_issue" "jira" {
  workflow_id = rootly_workflow_incident.jira.id
  task_params {
    title       = "{{ incident.title }}"
    description = "{{ incident.summary }}"
    project_key = "ROOT"
    issue_type = {
      id   = "10001"
      name = "Task"
    }
    status = {
      id   = "10000"
      name = "To Do"
    }
    labels = "{{ incident.environment_slugs | concat: incident.service_slugs | concat: incident.functionality_slugs | concat: incident.group_slugs | join: \",\" }}"
  }
}
