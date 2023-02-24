resource "rootly_workflow_incident" "my-workflow" {
  name        = "Trigger when an incident is created"
  description = "This workflow will trigger when an incident is created"
  trigger_params {
    triggers                  = ["incident_created"]
    incident_condition_kind   = "IS"
    incident_kinds            = ["normal"]
    incident_condition_status = "IS"
    incident_statuses         = ["started"]
  }
  enabled = true
}
