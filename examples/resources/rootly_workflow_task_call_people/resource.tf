resource "rootly_workflow_incident" "call_people" {
  name        = "Call a teammate"
  description = "Automatically call a specific teammate."
  trigger_params {
    triggers                  = ["incident_created"]
    incident_statuses         = ["started"]
    incident_condition_status = "IS"
  }
  enabled = true
}

resource "rootly_workflow_task_call_people" "call_people" {
  workflow_id = rootly_workflow_incident.call_people.id
  task_params {
    name    = "Call people"
    content = "We have an ongoing incident {{ incident.title }} of severity {{ incident.severity }} and your assistance is required."
  }
  # TODO
}
