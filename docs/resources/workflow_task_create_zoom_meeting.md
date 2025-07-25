---
page_title: "Resource rootly_workflow_task_create_zoom_meeting - terraform-provider-rootly"
subcategory: Workflow Tasks
description: |-
    Manages workflow create_zoom_meeting task.
---

# Resource (rootly_workflow_task_create_zoom_meeting)

Manages workflow create_zoom_meeting task.



## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import rootly_workflow_task_create_zoom_meeting using the `id`. For example:

```terraform
import {
  to = rootly_workflow_task_create_zoom_meeting.my-resource
  id = "00000000-0000-0000-0000-000000000000"
}
```

Locate the resource id in the web app, or retrieve it by listing resources through the API if it's not visible in the web app.

You can generate HCL from the import block using the `-generate-config-out` flag:

```console
terraform plan -generate-config-out=generated.tf
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

- `topic` (String) The meeting topic

Optional:

- `alternative_hosts` (List of String)
- `auto_recording` (String) Value must be one of `none`, `local`, `cloud`.
- `create_as_email` (String) The email to use if creating as email
- `password` (String) The meeting password
- `post_to_incident_timeline` (Boolean) Value must be one of true or false
- `post_to_slack_channels` (Block List) (see [below for nested schema](#nestedblock--task_params--post_to_slack_channels))
- `record_meeting` (Boolean) Rootly AI will record the meeting and automatically generate a transcript and summary from your meeting. Value must be one of true or false
- `task_type` (String)

<a id="nestedblock--task_params--post_to_slack_channels"></a>
### Nested Schema for `task_params.post_to_slack_channels`

Required:

- `id` (String)
- `name` (String)
