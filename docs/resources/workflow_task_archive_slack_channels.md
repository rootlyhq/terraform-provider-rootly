---
page_title: "Resource rootly_workflow_task_archive_slack_channels - terraform-provider-rootly"
subcategory: Workflow Tasks
description: |-
    Manages workflow archive_slack_channels task.
---

# Resource (rootly_workflow_task_archive_slack_channels)

Manages workflow archive_slack_channels task.

## Example Usage

```shell
resource "rootly_workflow_incident" "auto_archive_incident" {
  name        = "Auto archive incident 48hrs after resolution"
  description = "After an incident has been resolved, automatically archive the channel after 48 hours."
  trigger_params {
    triggers                  = ["status_updated"]
    wait                      = "2 days"
    incident_statuses         = ["resolved"]
    incident_condition_status = "IS"
  }
  enabled = true
}

resource "rootly_workflow_task_archive_slack_channels" "archive_slack_channels" {
  workflow_id     = rootly_workflow_incident.auto_archive_incident.id
  skip_on_failure = false
  enabled         = true
  task_params {
    name = "Archive Slack channels"
    channels {
      id   = "{{ incident.slack_channel_id }}"
      name = "{{ incident.slack_channel_id }}"
    }
  }
}
```

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import rootly_workflow_task_archive_slack_channels using the `id`. For example:

```terraform
import {
  to = rootly_workflow_task_archive_slack_channels.my-resource
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

- `channels` (Block List, Min: 1) (see [below for nested schema](#nestedblock--task_params--channels))

Optional:

- `task_type` (String)

<a id="nestedblock--task_params--channels"></a>
### Nested Schema for `task_params.channels`

Required:

- `id` (String)
- `name` (String)
