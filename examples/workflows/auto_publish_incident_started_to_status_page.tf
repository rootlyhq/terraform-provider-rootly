resource "rootly_workflow_incident" "auto_publish_incident_started_to_status_page" {
  name        = "Auto Update to Status Page - Incident Started"
  description = "Automatically updates your status page once an incident starts"
  trigger_params {
    triggers                  = ["incident_created"]
    incident_statuses         = ["started"]
    incident_condition_status = "IS"
  }
  enabled = true
}

resource "rootly_workflow_task_publish_incident" "publish_incident" {
  workflow_id = rootly_workflow_incident.auto_publish_incident_started_to_status_page.id
  task_params {
    name = "Publish incident started to status page"
    incident = {
      id   = "{{ incident.id }}"
      name = "{{ incident.id }}"
    }
    public_title = "{{ incident.title }}"
    event        = <<EOT
The incident {{ incident.title }} is ongoing and being investigated. 

Severity: {{ incident.severity }}
Summary: {{ incident.summary }}
EOT
    status       = "investigating"

  }
}
