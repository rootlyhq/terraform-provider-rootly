resource "rootly_workflow_incident" "canvas" {
  name    = "Slack Canvas report"
  enabled = false

  trigger_params {
    triggers = ["slack_channel_created"]
  }
}

resource "rootly_workflow_task_create_slack_canvas" "report" {
  workflow_id = rootly_workflow_incident.canvas.id
  name        = "Create incident Canvas"
  enabled     = false
  position    = 1

  task_params {
    title = "Incident {{ incident.sequential_id }}: {{ incident.title }}"
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
