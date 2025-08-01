---
page_title: "Resource rootly_workflow_task_create_google_calendar_event - terraform-provider-rootly"
subcategory: Workflow Tasks
description: |-
    Manages workflow create_google_calendar_event task.
---

# Resource (rootly_workflow_task_create_google_calendar_event)

Manages workflow create_google_calendar_event task.

## Example Usage

```shell
resource "rootly_workflow_incident" "schedule_postmortem_review_meeting" {
  name        = "Schedule Postmortem Review Meeting"
  description = "Automatically schedule a Google Calendar meeting to review the postmortem."
  trigger_params {
    triggers                  = ["status_updated"]
    incident_statuses         = ["resolved"]
    incident_condition_status = "IS"
  }
  enabled = true
}

resource "rootly_workflow_task_create_google_calendar_event" "create_google_calendar_event" {
  workflow_id     = rootly_workflow_incident.schedule_postmortem_review_meeting.id
  skip_on_failure = false
  enabled         = true
  task_params {
    name               = "Schedule Postmortem Review Meeting"
    days_until_meeting = 7
    meeting_duration   = "60min"
    summary            = "#{{ incident.sequential_id }} {{ incident.title }} Postmortem Review"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `task_params` (Block List, Min: 1, Max: 1) The parameters for this workflow task. (see [below for nested schema](#nestedblock--task_params))
- `workflow_id` (String) The ID of the parent workflow

### Optional

- `enabled` (Boolean) Enable/disable this workflow task
- `name` (String) Name of the workflow task
- `position` (Number) The position of the workflow task (1 being top of list)
- `skip_on_failure` (Boolean) Skip workflow task if any failures

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--task_params"></a>
### Nested Schema for `task_params`

Required:

- `days_until_meeting` (String) The days until meeting
- `description` (String) The event description
- `meeting_duration` (String) Meeting duration in format like '1 hour', '30 minutes'
- `summary` (String) The event summary
- `time_of_meeting` (String) Time of meeting in format HH:MM

Optional:

- `attendees` (List of String) Emails of attendees
- `calendar_id` (String)
- `can_guests_invite_others` (Boolean) Value must be one of true or false
- `can_guests_modify_event` (Boolean) Value must be one of true or false
- `can_guests_see_other_guests` (Boolean) Value must be one of true or false
- `conference_solution_key` (String) Sets the video conference type attached to the meeting. Value must be one of `eventHangout`, `eventNamedHangout`, `hangoutsMeet`, `addOn`.
- `exclude_weekends` (Boolean) Value must be one of true or false
- `post_to_incident_timeline` (Boolean) Value must be one of true or false
- `post_to_slack_channels` (Block List) (see [below for nested schema](#nestedblock--task_params--post_to_slack_channels))
- `send_updates` (Boolean) Send an email to the attendees notifying them of the event. Value must be one of true or false
- `task_type` (String)
- `time_zone` (String) A valid IANA time zone name.

<a id="nestedblock--task_params--post_to_slack_channels"></a>
### Nested Schema for `task_params.post_to_slack_channels`

Required:

- `id` (String)
- `name` (String)

## Import

rootly_workflow_task_create_google_calendar_event can be imported using the [`import` command](https://developer.hashicorp.com/terraform/cli/commands/import).

```sh
terraform import rootly_workflow_task_create_google_calendar_event.primary a816421c-6ceb-481a-87c4-585e47451f24
```

Or using an [`import` block](https://developer.hashicorp.com/terraform/language/import).

```terraform
import {
  to = rootly_workflow_task_create_google_calendar_event.primary
  id = "a816421c-6ceb-481a-87c4-585e47451f24"
}
```

Locate the resource id in the web app, or retrieve it by listing resources through the API if it's not visible in the web app.

HCL can be generated from the import block using the `-generate-config-out` flag.

```sh
terraform plan -generate-config-out=generated.tf
```
