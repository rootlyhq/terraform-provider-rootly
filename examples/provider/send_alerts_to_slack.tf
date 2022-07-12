resource "rootly_workflow_alert" "send_alerts_to_slack" {
  name = "Notify alerts channel"
  description = "Sends a customizable message block to Slack with details of your alerts"
  trigger_params {
    triggers = ["alert_created"]
    
  }
  enabled = true
}

resource "rootly_workflow_task_send_slack_blocks" "send_slack_blocks" {
  workflow_id = rootly_workflow_alert.send_alerts_to_slack.id
  task_params {
    name = "Send alert block"
    message = ":boom: New alert!"
    blocks {
      id = "undefined"
      name = "undefined"
    }
    blocks {
      id = "undefined"
      name = "undefined"
    }
  }
}