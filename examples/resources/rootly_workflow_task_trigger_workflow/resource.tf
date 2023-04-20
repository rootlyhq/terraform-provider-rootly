resource "rootly_workflow_incident" "trigger_another_workflow" {
  name        = "Trigger another workflow"
  description = "Trigger another workflow"
  trigger_params {
    triggers                  = ["incident_created"]
    incident_statuses         = ["started"]
    incident_condition_status = "IS"
  }
  enabled = true
}

resource "rootly_workflow_task_trigger_workflow" "trigger_another_workflow_task" {
  position    = 1
  workflow_id = rootly_workflow_incident.trigger_another_workflow.id

  task_params {
    kind = "incident"
    workflow = {
      "id"   = rootly_workflow_incident.another_workflow.id
      "name" = "Trigger another workflow"
    }
    resource = {
      "id"   = "{{ incident.id }}"
      "name" = "{{ incident.id }}"
    }
  }
}