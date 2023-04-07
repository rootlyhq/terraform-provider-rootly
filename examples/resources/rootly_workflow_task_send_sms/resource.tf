resource "rootly_workflow_incident" "send_sms" {
  name        = "Send SMS to a teammate"
  description = "Automatically send SMS to a specific teammate."
  trigger_params {
    triggers                  = ["incident_created"]
    incident_statuses         = ["started"]
    incident_condition_status = "IS"
  }
  enabled = true
}

resource "rootly_workflow_task_send_sms" "send_sms" {
  workflow_id     = rootly_workflow_incident.send_sms.id
  skip_on_failure = false
  enabled         = true
  task_params {
    name          = "Send SMS"
    content       = "We have an ongoing incident {{ incident.title }} of severity {{ incident.severity }} and your assistance is required."
    phone_numbers = ["+11231231231", "+11231231232"]
  }
}
