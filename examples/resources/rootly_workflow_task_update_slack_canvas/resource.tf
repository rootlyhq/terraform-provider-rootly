resource "rootly_workflow_incident" "canvas" {
  name    = "Slack Canvas report"
  enabled = false

  trigger_params {
    triggers = ["incident_updated"]
  }
}

resource "rootly_workflow_task_update_slack_canvas" "report" {
  workflow_id = rootly_workflow_incident.canvas.id
  name        = "Update incident Canvas"
  enabled     = false
  position    = 1

  task_params {
    operation = "insert_at_end"
    channel {
      id   = "{{ incident.slack_channel_id }}"
      name = "Incident channel"

      workspace {
        id   = "T0123456789"
        name = "Engineering"
      }
    }
    content = <<-MARKDOWN
      # Incident update

      **Summary:** {{ incident.summary }}
    MARKDOWN
  }
}
