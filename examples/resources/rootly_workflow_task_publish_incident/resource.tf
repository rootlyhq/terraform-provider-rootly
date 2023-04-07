resource "rootly_workflow_incident" "auto_publish_incident_resolved_to_status_page" {
  name        = "Auto Update to Status Page - Incident Resolved"
  description = "Automatically updates your status page once an incident is resolved"
  trigger_params {
    triggers                  = ["status_updated"]
    incident_statuses         = ["resolved"]
    incident_condition_status = "IS"
  }
  enabled = true
}

resource "rootly_workflow_task_publish_incident" "publish_incident" {
  workflow_id     = rootly_workflow_incident.auto_publish_incident_resolved_to_status_page.id
  skip_on_failure = false
  enabled         = true
  task_params {
    name = "Publish incident resolved to status page"
    incident = {
      id   = "{{ incident.id }}"
      name = "{{ incident.id }}"
    }
    public_title = "{{ incident.title }}"
    event        = <<EOT
The incident {{ incident.title }} is resolved and postmortem being drafted. 

Severity: {{ incident.severity }}
Summary: {{ incident.summary }}
EOT
    status       = "resolved"

  }
}
