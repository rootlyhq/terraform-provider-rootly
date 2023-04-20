resource "rootly_workflow_pulse" "my-workflow" {
  name        = "Trigger when a pulse is created"
  description = "This workflow will trigger when a pulse is created"
  trigger_params {
    triggers = ["pulse_created"]

  }
  enabled = true
}
