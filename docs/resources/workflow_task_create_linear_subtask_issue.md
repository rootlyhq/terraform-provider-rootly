---
page_title: "Resource rootly_workflow_task_create_linear_subtask_issue - terraform-provider-rootly"
subcategory: ""
description: |-
    Manages workflow createlinearsubtask_issue task.
---

# Resource (rootly_workflow_task_create_linear_subtask_issue)

Manages workflow create_linear_subtask_issue task.

## Example Usage

```
resource "rootly_workflow_task_create_linear_subtask_issue" "foo" {
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

- `parent_issue_id` (String) The parent issue.
- `state` (Map of String) The state id and display name.
- `title` (String) The issue title.

Optional:

- `description` (String) The issue description.
- `priority` (Map of String) The priority id and display name.
- `task_type` (String)

## Import

Import is supported using the following syntax:

```shell
terraform import rootly_workflow_task_create_linear_subtask_issue.foo 11111111-2222-3333-4444-555555555555
```