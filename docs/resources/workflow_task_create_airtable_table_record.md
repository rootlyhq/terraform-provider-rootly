---
page_title: "Resource rootly_workflow_task_create_airtable_table_record - terraform-provider-rootly"
subcategory: ""
description: |-
    Manages workflow createairtabletable_record task.
---

# Resource (rootly_workflow_task_create_airtable_table_record)

Manages workflow create_airtable_table_record task.

## Example Usage

```
resource "rootly_workflow_task_create_airtable_table_record" "foo" {
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

- `base_key` (String) The base key
- `table_name` (String) The table name

Optional:

- `custom_fields_mapping` (String) Custom field mappings. Can contain liquid markup and need to be valid JSON.
- `task_type` (String)

## Import

Import is supported using the following syntax:

```shell
terraform import rootly_workflow_task_create_airtable_table_record.foo 11111111-2222-3333-4444-555555555555
```