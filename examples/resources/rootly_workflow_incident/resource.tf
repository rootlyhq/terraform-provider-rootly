data "rootly_severity" "critical" {
  slug = "sev0"
}

resource "rootly_workflow_incident" "my-workflow" {
  name        = "Trigger when an incident is created and severity is critical"
  description = "This workflow will trigger when an incident is created and severity is critical"
  trigger_params {
    triggers                    = ["incident_created"]
    incident_condition_kind     = "IS"
    incident_kinds              = ["normal"]
    incident_condition_status   = "IS"
    incident_statuses           = ["started"]
    incident_condition_severity = "IS"
  }
  severity_ids = [data.rootly_severity.critical.id]
  enabled      = true
}

# Route failure notifications for this workflow to the incident's Slack channel
# and a dedicated ops channel instead of the account-wide default channel.
resource "rootly_workflow_incident" "notify-on-failure" {
  name = "Page on-call when incident is created"
  trigger_params {
    triggers = ["incident_created"]
  }

  failure_notification_mode = "custom"
  failure_notification_channels {
    id   = "{{ incident.slack_channel_id }}"
    name = "{{ incident.slack_channel_id }}"
  }
  failure_notification_channels {
    id   = "C0123456789"
    name = "workflow-failures"
  }
}

# Suppress failure notifications for a noisy workflow.
resource "rootly_workflow_incident" "quiet" {
  name = "Best-effort enrichment"
  trigger_params {
    triggers = ["incident_updated"]
  }

  failure_notification_mode = "off"
}
