---
page_title: "Resource rootly_workflow_task_create_incident_postmortem - terraform-provider-rootly"
subcategory: Workflow Tasks
description: |-
    Manages workflow createincidentpostmortem task.
---

# Resource (rootly_workflow_task_create_incident_postmortem)

Manages workflow create_incident_postmortem task.



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

- `incident_id` (String) UUID of the incident that needs a retrospective.
- `title` (String) The retrospective title

Optional:

- `status` (String)
- `task_type` (String)
- `template` (Map of String) Map must contain two fields, `id` and `name`. Retrospective template to use.