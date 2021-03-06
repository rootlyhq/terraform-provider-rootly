---
page_title: "Resource rootly_workflow_task_add_action_item - terraform-provider-rootly"
subcategory: ""
description: |-
    Manages workflow addactionitem task.
---

# Resource (rootly_workflow_task_add_action_item)

Manages workflow add_action_item task.

## Example Usage

```
resource "rootly_workflow_task_add_action_item" "foo" {
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

- `priority` (String) The action item priority.
- `status` (String) The action item status.
- `summary` (String) The action item summary.

Optional:

- `assigned_to_user_id` (String) The user id this action item is assigned to
- `description` (String) The action item description.
- `kind` (String) The action item kind.
- `post_to_incident_timeline` (Boolean)
- `post_to_slack_channels` (Block List) (see [below for nested schema](#nestedblock--task_params--post_to_slack_channels))
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
terraform import rootly_workflow_task_add_action_item.foo 11111111-2222-3333-4444-555555555555
```