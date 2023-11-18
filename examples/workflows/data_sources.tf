# uses output of severity data source as input for workflow
data "rootly_severity" "critical" {
  slug = "sev0"
}

resource "rootly_workflow_incident" "ping_oncall" {
  name        = "Ping on-call when critical incident"
  description = "ping on-call when critical incident happens"
  trigger_params {
    triggers                  = ["incident_created"]
    incident_condition_kind   = "IS"
    incident_kinds            = ["normal"]
    incident_condition_status = "IS"
    incident_statuses         = ["started"]
    severity_ids              = [data.rootly_severity.critical.id]
  }
  enabled = true
}

resource "rootly_workflow_task_send_sms" "sms_oncall" {
  workflow_id = rootly_workflow_incident.ping_oncall.id
  name        = "On-call team"

  task_params {
    phone_numbers = ["+11231231234"]
    content       = "Critical incident started"
  }
}
