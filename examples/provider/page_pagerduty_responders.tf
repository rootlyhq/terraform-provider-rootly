resource "rootly_workflow_incident" "page_pagerduty_responders" {
  name = "Page responders via PagerDuty"
  description = "Automatically page responders to join the incident depending on what's been impacted (see conditions)."
  trigger_params {
    triggers = ["incident_created"]
    incident_statuses = ["started"]
    incident_condition_status = "IS"
  }
  enabled = true
}

resource "rootly_workflow_task_page_pagerduty_on_call_responders" "page_pagerduty_on_call_responders" {
  workflow_id = rootly_workflow_incident.page_pagerduty_responders.id
  task_params {
    name = "Page PagerDuty on-call responders"
  }
}