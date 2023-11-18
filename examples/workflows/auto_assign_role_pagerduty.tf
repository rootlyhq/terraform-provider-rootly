resource "rootly_workflow_incident" "auto_assign_role_pagerduty" {
  name        = "Auto Assign Roles from On-call Rotation"
  description = "Automatically assign users to specific roles (e.g. Commander) based on your on-call rotation."
  trigger_params {
    triggers                  = ["incident_created"]
    incident_statuses         = ["started"]
    incident_condition_status = "IS"
  }
  enabled = true
}

resource "rootly_workflow_task_auto_assign_role_pagerduty" "auto_assign_role_pagerduty" {
  workflow_id = rootly_workflow_incident.auto_assign_role_pagerduty.id
  name        = "Auto Assign Roles from On-call Rotation"

  task_params {
  }
}
