---
page_title: "Resource rootly_workflow_task_update_incident_postmortem - terraform-provider-rootly"
subcategory: Workflow Tasks
description: |-
    Manages workflow update_incident_postmortem task.
---

# Resource (rootly_workflow_task_update_incident_postmortem)

Manages workflow update_incident_postmortem task.



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

- `postmortem_id` (String) UUID of the retrospective that needs to be updated

Optional:

- `status` (String)
- `task_type` (String)
- `title` (String) The incident title

## Import

rootly_workflow_task_update_incident_postmortem can be imported using the [`import` command](https://developer.hashicorp.com/terraform/cli/commands/import).

```sh
terraform import rootly_workflow_task_update_incident_postmortem.primary a816421c-6ceb-481a-87c4-585e47451f24
```

Or using an [`import` block](https://developer.hashicorp.com/terraform/language/import).

```terraform
import {
  to = rootly_workflow_task_update_incident_postmortem.primary
  id = "a816421c-6ceb-481a-87c4-585e47451f24"
}
```

Locate the resource id in the web app, or retrieve it by listing resources through the API if it's not visible in the web app.

HCL can be generated from the import block using the `-generate-config-out` flag.

```sh
terraform plan -generate-config-out=generated.tf
```
