---
page_title: "Resource rootly_workflow_task_invite_to_slack_channel_victor_ops - terraform-provider-rootly"
subcategory: ""
description: |-
    Manages workflow invitetoslackchannelvictor_ops task.
---

# Resource (rootly_workflow_task_invite_to_slack_channel_victor_ops)

Manages workflow invite_to_slack_channel_victor_ops task.

## Example Usage

```
resource "rootly_workflow_task_invite_to_slack_channel_victor_ops" "foo" {
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

- `schedule` (Map of String)

Optional:

- `channels` (List of String)
- `task_type` (String)

## Import

Import is supported using the following syntax:

```shell
terraform import rootly_workflow_task_invite_to_slack_channel_victor_ops.foo 11111111-2222-3333-4444-555555555555
```