---
page_title: "Resource rootly_workflow_task_snapshot_looker_look - terraform-provider-rootly"
subcategory: ""
description: |-
    Manages workflow snapshotlookerlook task.
---

# Resource (rootly_workflow_task_snapshot_looker_look)

Manages workflow snapshot_looker_look task.

## Example Usage

```
resource "rootly_workflow_task_snapshot_looker_look" "foo" {
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

- `dashboards` (Block List, Min: 1) (see [below for nested schema](#nestedblock--task_params--dashboards))

Optional:

- `post_to_slack_channels` (Block List) (see [below for nested schema](#nestedblock--task_params--post_to_slack_channels))
- `task_type` (String)

<a id="nestedblock--task_params--dashboards"></a>
### Nested Schema for `task_params.dashboards`

Required:

- `name` (String)

Read-Only:

- `id` (String) The ID of this resource.


<a id="nestedblock--task_params--post_to_slack_channels"></a>
### Nested Schema for `task_params.post_to_slack_channels`

Required:

- `name` (String)

Read-Only:

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import rootly_workflow_task_snapshot_looker_look.foo 11111111-2222-3333-4444-555555555555
```