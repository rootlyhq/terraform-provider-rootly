---
page_title: "Resource rootly_workflow_task_invite_to_slack_channel - terraform-provider-rootly"
subcategory: ""
description: |-
    Manages workflow invitetoslack_channel task.
---

# Resource (rootly_workflow_task_invite_to_slack_channel)

Manages workflow invite_to_slack_channel task.

## Example Usage

```
resource "rootly_workflow_task_invite_to_slack_channel" "foo" {
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

- `channel` (Map of String)

Optional:

- `slack_user_groups` (Block List) (see [below for nested schema](#nestedblock--task_params--slack_user_groups))
- `slack_users` (Block List) (see [below for nested schema](#nestedblock--task_params--slack_users))
- `task_type` (String)

<a id="nestedblock--task_params--slack_user_groups"></a>
### Nested Schema for `task_params.slack_user_groups`

Required:

- `name` (String)

Read-Only:

- `id` (String) The ID of this resource.


<a id="nestedblock--task_params--slack_users"></a>
### Nested Schema for `task_params.slack_users`

Required:

- `name` (String)

Read-Only:

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import rootly_workflow_task_invite_to_slack_channel.foo 11111111-2222-3333-4444-555555555555
```