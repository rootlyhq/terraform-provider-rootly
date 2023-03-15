data "rootly_severity" "critical" {
  slug = "sev0"
}

resource "rootly_workflow_incident" "my-workflow" {
  name        = "Trigger when an incident is created and severity is critical"
  description = "This workflow will trigger when an incident is created and severity is critical"
  trigger_params {
    triggers                    = ["incident_created"]
    incident_condition_kind     = "IS"
    incident_kinds              = ["normal"]
    incident_condition_status   = "IS"
    incident_statuses           = ["started"]
    incident_condition_severity = "IS"
  }
  severity_ids = [data.rootly_severity.critical.id]
  enabled      = true
}
