resource "rootly_workflow_incident" "page_rootly_on_call_responders" {
  name        = "Page Rootly on-call responders"
  description = "Automatically page Rootly on-call responders when an incident is created."

  trigger_params {
    triggers                  = ["incident_created"]
    incident_statuses         = ["started"]
    incident_condition_status = "IS"
  }

  enabled = true
}

resource "rootly_workflow_task_page_rootly_on_call_responders" "page_rootly_on_call_responders" {
  workflow_id     = rootly_workflow_incident.page_rootly_on_call_responders.id
  skip_on_failure = true
  enabled         = true

  task_params {
    alert_urgency_id = "67f9be75-3b5e-40be-8fd8-4e27d6a7f75d"
    summary          = "Incident: {{ incident.severity }} - {{ incident.title }}"
    description      = "{{ incident.summary }}"

    group_target = {
      id   = "5ec9d015-20bf-440f-932a-d46b50fd6430"
      name = "Customer Relations"
    }
  }
}
