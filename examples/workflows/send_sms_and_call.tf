resource "rootly_workflow_incident" "send_sms_and_call" {
  name        = "Send SMS and call a teammate"
  description = "Automatically call and send SMS to a specific teammate."
  trigger_params {
    triggers                  = ["incident_created"]
    incident_statuses         = ["started"]
    incident_condition_status = "IS"
  }
  enabled = true
}

resource "rootly_workflow_task_send_sms" "send_sms" {
  workflow_id = rootly_workflow_incident.send_sms_and_call.id
  task_params {
    name    = "Send SMS"
    content = "We have an ongoing incident {{ incident.title }} of severity {{ incident.severity }} and your assistance is required."
  }
}

resource "rootly_workflow_task_call_people" "call_people" {
  workflow_id = rootly_workflow_incident.send_sms_and_call.id
  task_params {
    name    = "Call people"
    content = "We have an ongoing incident {{ incident.title }} of severity {{ incident.severity }} and your assistance is required."
  }
}
