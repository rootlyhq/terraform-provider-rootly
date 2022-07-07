---
page_title: "Resource rootly_workflow_task_create_shortcut_task - terraform-provider-rootly"
subcategory: ""
description: |-
    Manages workflow createshortcuttask task.
---

# Resource (rootly_workflow_task_create_shortcut_task)

Manages workflow create_shortcut_task task.

## Example Usage

```
resource "rootly_workflow_task_create_shortcut_task" "foo" {
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

- `completion` (Map of String) The completion id and display name.
- `description` (String) The task description
- `parent_story_id` (String) The parent story

Optional:

- `task_type` (String)

## Import

Import is supported using the following syntax:

```shell
terraform import rootly_workflow_task_create_shortcut_task.foo 11111111-2222-3333-4444-555555555555
```