resource "rootly_workflow_incident" "schedule_postmortem_review_meeting" {
  name = "Schedule Postmortem Review Meeting"
  description = "Automatically schedule a Google Calendar meeting to review the postmortem."
  trigger_params {
    triggers = ["status_updated"]
    incident_statuses = ["resolved"]
    incident_condition_status = "IS"
  }
  enabled = true
}

resource "rootly_workflow_task_create_google_calendar_event" "create_google_calendar_event" {
  workflow_id = rootly_workflow_incident.schedule_postmortem_review_meeting.id
  task_params {
    name = "Schedule Postmortem Review Meeting"
    days_until_meeting = 7
    meeting_duration = "60min"
    summary = "#{{ incident.sequential_id }} {{ incident.title }} Postmortem Review"
  }
}