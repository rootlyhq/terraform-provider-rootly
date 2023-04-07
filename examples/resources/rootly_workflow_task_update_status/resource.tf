resource "rootly_workflow_incident" "auto_resolve_incident" {
  name        = "Auto resolve inactive incidents"
  description = "After an incident has been inactive for 24 hrs, automatically mark it as resolved."
  trigger_params {
    triggers                  = ["incident_created"]
    wait                      = "24 hours"
    incident_statuses         = ["resolved", "mitigated"]
    incident_condition_status = "ANY"
  }
  enabled = true
}

resource "rootly_workflow_task_update_status" "update_status" {
  workflow_id     = rootly_workflow_incident.auto_resolve_incident.id
  skip_on_failure = false
  enabled         = true
  task_params {
    status  = "resolved"
    message = "Automatically marked as resolved due to inactivity"
  }
}
