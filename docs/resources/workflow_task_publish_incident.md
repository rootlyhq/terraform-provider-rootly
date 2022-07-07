---
page_title: "Resource rootly_workflow_task_publish_incident - terraform-provider-rootly"
subcategory: ""
description: |-
    Manages workflow publish_incident task.
---

# Resource (rootly_workflow_task_publish_incident)

Manages workflow publish_incident task.

## Example Usage

```
resource "rootly_workflow_task_publish_incident" "foo" {
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

- `event` (String) Incident event description
- `incident` (Map of String)
- `public_title` (String)
- `status` (String)
- `status_page_ids` (List of String)

Optional:

- `task_type` (String)

## Import

Import is supported using the following syntax:

```shell
terraform import rootly_workflow_task_publish_incident.foo 11111111-2222-3333-4444-555555555555
```