resource "rootly_workflow_incident" "publish_incident_to_status_page" {
  name        = "Publish incident to status page"
  description = "Automatically publish your incident to one or many status pages."
  trigger_params {
    triggers                  = ["incident_created"]
    incident_statuses         = ["started"]
    incident_condition_status = "IS"
  }
  enabled = true
}

resource "rootly_workflow_task_publish_incident" "publish_incident" {
  workflow_id = rootly_workflow_incident.publish_incident_to_status_page.id
  name        = "Automatically publish to status pages"

  task_params {
  }
}
