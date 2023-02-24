resource "rootly_workflow_alert" "my-workflow" {
  name        = "Trigger when an alert is created"
  description = "This workflow will trigger when an alert is created"
  trigger_params {
    triggers = ["alert_created"]

  }
  enabled = true
}
