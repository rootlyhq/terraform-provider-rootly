---
page_title: "Resource rootly_workflow_task_update_incident_postmortem - terraform-provider-rootly"
subcategory: Workflow Tasks
description: |-
    Manages workflow updateincidentpostmortem task.
---

# Resource (rootly_workflow_task_update_incident_postmortem)

Manages workflow update_incident_postmortem task.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `task_params` (Block List, Min: 1, Max: 1) The parameters for this workflow task. (see [below for nested schema](#nestedblock--task_params))
- `workflow_id` (String) The ID of the parent workflow

### Optional

- `position` (Number) The position of the workflow task (1 being top of list)

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--task_params"></a>
### Nested Schema for `task_params`

Optional:

- `postmortem_id` (String) UUID of the postmortem that needs to be updated.
- `status` (String)
- `task_type` (String)
- `title` (String) The incident title