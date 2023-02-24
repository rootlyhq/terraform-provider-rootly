resource "rootly_workflow_incident" "send_email_to_stakeholders" {
  name        = "Send update email"
  description = "Workflow for sending an email to stakeholders (e.g. leadership, legal) to keep them updated on the incident."
  trigger_params {
    triggers                  = ["incident_created"]
    wait                      = "30 seconds"
    incident_statuses         = ["started"]
    incident_condition_status = "IS"
  }
  enabled = true
}

resource "rootly_workflow_task_send_email" "send_email" {
  workflow_id = rootly_workflow_incident.send_email_to_stakeholders.id
  task_params {
    name      = "Send Email to Subscribers"
    to        = ["{{ incident.subscribers | map: 'email' | join: ',' }}", "{{ incident.raw_severity | get: 'notify_emails' | join: ',' }}", "{{ incident.raw_environments | map: 'notify_emails' | flatten | join: ',' }}", "{{ incident.raw_functionalities | map: 'notify_emails' | flatten | join: ',' }}", "{{ incident.raw_services | map: 'notify_emails' | flatten | join: ',' }}", "{{ incident.raw_types | map: 'notify_emails' | flatten | join: ',' }}", "{{ incident.raw_groups | map: 'notify_emails' | flatten | join: ',' }}"]
    preheader = "{{ incident.summary }}"
    subject   = "#{{ incident.sequential_id }} {{ incident.title }} status changed to: {{ incident.status }}"
    body      = <<EOT
## #{{ incident.sequential_id }} {{ incident.title }} status changed to: {{ incident.status }}

{% if incident.status == "mitigated" %}
  {% if incident.mitigation_message != blank %}
    How? {{ incident.mitigation_message }}
  {% endif %}
{% elsif incident.status == "resolved" %}
  {% if incident.resolution_message != blank %}
    How? {{ incident.resolution_message }}
  {% endif %}
{% elsif incident.status == "cancelled" %}
  {% if incident.cancellation_message != blank %}
    How? {{ incident.cancellation_message }}
  {% endif %}
{% endif %}

You can use the following link to see the incident.

{% if incident.url != blank %}
  [View Incident]({{ incident.url }} 'btn')
{% endif %}

{% if incident.private_status_page_url != blank %}
  [Rootly Private Status Page]({{ incident.private_status_page_url }})
{% endif %}

{% if incident.slack_channel_url != blank %}
  ![Slack!]({{ images.slack_url_16 }}) [#{{ incident.slack_channel_name }}]({{ incident.slack_channel_url }})
{% endif %}

{% if incident.google_meeting_url != blank %}
  ![Google Meet!]({{ images.google_meet_url_16 }}) [Google Meet Room]({{ incident.google_meeting_url }})
{% endif %}

{% if incident.zoom_meeting_join_url != blank %}
  ![Zoom Room!]({{ images.zoom_url_16 }}) [Zoom Room]({{ incident.zoom_meeting_join_url }})
{% endif %}

{% if incident.webex_meeting_url != blank %}
  ![Webex Meeting!]({{ images.webex_url_16 }}) [Webex Meeting]({{ incident.webex_meeting_url }})
{% endif %}

{% if incident.jira_issue_url != blank %}
  ![Jira Issue!]({{ images.jira_url_16 }}) [Jira Issue]({{ incident.jira_issue_url }})
{% endif %}

{% if incident.asana_task_url != blank %}
  ![Asana Task!]({{ images.asana_url_16 }}) [Asana Task]({{ incident.asana_task_url }})
{% endif %}

{% if incident.github_issue_url != blank %}
  ![Asana Task!]({{ images.asana_url_16 }}) [Github issue]({{ incident.github_issue_url }})
{% endif %}

{% if incident.trello_card_url != blank %}
  ![Trello card!]({{ images.trello_url_16 }}) [Trello card]({{ incident.trello_card_url }})
{% endif %}

{% if incident.shortcut_story_url != blank %}
  ![Shortcut story!]({{ images.shortcut_url_16 }}) [Shortcut story]({{ incident.shortcut_story_url }})
{% endif %}

{% if incident.service_now_incident_url != blank %}
  ![ServiceNow Incident!]({{ images.service_now_url_16 }}) [ServiceNow Incident]({{ incident.service_now_incident_url }})
{% endif %}

If you have any questions, please don't hesitate to [send us an email](mailto:support@rootly.com).

EOT
  }
}
