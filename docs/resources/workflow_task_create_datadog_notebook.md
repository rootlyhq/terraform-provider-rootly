---
page_title: "Resource rootly_workflow_task_create_datadog_notebook - terraform-provider-rootly"
subcategory: ""
description: |-
    Manages workflow createdatadognotebook task.
---

# Resource (rootly_workflow_task_create_datadog_notebook)

Manages workflow create_datadog_notebook task.

## Example Usage

```
resource "rootly_workflow_task_create_datadog_notebook" "foo" {
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

- `kind` (String) The notebook kind
- `title` (String) The notebook title

Optional:

- `content` (String) The notebook content
- `post_mortem_template_id` (String) Post mortem template to use when creating notebook, if desired.
- `task_type` (String)

## Import

Import is supported using the following syntax:

```shell
terraform import rootly_workflow_task_create_datadog_notebook.foo 11111111-2222-3333-4444-555555555555
```