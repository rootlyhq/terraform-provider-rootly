---
page_title: "Resource rootly_workflow_task_get_pulses - terraform-provider-rootly"
subcategory: ""
description: |-
    Manages workflow get_pulses task.
---

# Resource (rootly_workflow_task_get_pulses)

Manages workflow get_pulses task.

## Example Usage

```
resource "rootly_workflow_task_get_pulses" "foo" {
  name = "bar"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `task_params` (Block List, Min: 1, Max: 1) The parameters for this workflow task. (see [below for nested schema](#nestedblock--task_params))
- `workflow_id` (String) The ID of the parent workflow

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--task_params"></a>
### Nested Schema for `task_params`

Required:

- `past_duration` (String) How far back to fetch commits (in format '1 minute', '30 days', '3 months', etc.)

Optional:

- `environment_ids` (List of String)
- `labels` (List of String)
- `post_to_slack_channels` (Block List) (see [below for nested schema](#nestedblock--task_params--post_to_slack_channels))
- `refs` (List of String)
- `service_ids` (List of String)
- `sources` (List of String)
- `task_type` (String)

<a id="nestedblock--task_params--post_to_slack_channels"></a>
### Nested Schema for `task_params.post_to_slack_channels`

Required:

- `name` (String)

Read-Only:

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import rootly_workflow_task_get_pulses.foo 11111111-2222-3333-4444-555555555555
```