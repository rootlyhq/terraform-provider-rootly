---
page_title: "Resource rootly_workflow_task_create_service_now_incident - terraform-provider-rootly"
subcategory: ""
description: |-
    Manages workflow createservicenow_incident task.
---

# Resource (rootly_workflow_task_create_service_now_incident)

Manages workflow create_service_now_incident task.

## Example Usage

```
resource "rootly_workflow_task_create_service_now_incident" "foo" {
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

- `title` (String) The incident title

Optional:

- `completion` (Map of String) The completion id and display name.
- `custom_fields_mapping` (String) Custom field mappings. Can contain liquid markup and need to be valid JSON.
- `description` (String) The incident description
- `priority` (Map of String) The priority id and display name.
- `task_type` (String)

## Import

Import is supported using the following syntax:

```shell
terraform import rootly_workflow_task_create_service_now_incident.foo 11111111-2222-3333-4444-555555555555
```