---
page_title: "Resource rootly_workflow_task_update_gitlab_issue - terraform-provider-rootly"
subcategory: Workflow Tasks
description: |-
    Manages workflow updategitlabissue task.
---

# Resource (rootly_workflow_task_update_gitlab_issue)

Manages workflow update_gitlab_issue task.



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

- `completion` (Map of String) Map must contain two fields, `id` and `name`.
- `issue_id` (String) The issue id

Optional:

- `description` (String) The issue description
- `due_date` (String) The due date
- `issue_type` (String) The issue type. Value must be one of `issue`, `incident`, `test_case`, `task`.
- `labels` (String) The issue labels
- `task_type` (String)
- `title` (String) The issue title