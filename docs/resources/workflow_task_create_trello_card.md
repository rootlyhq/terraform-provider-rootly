---
page_title: "Resource rootly_workflow_task_create_trello_card - terraform-provider-rootly"
subcategory: Workflow Tasks
description: |-
    Manages workflow create_trello_card task.
---

# Resource (rootly_workflow_task_create_trello_card)

Manages workflow create_trello_card task.



## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import rootly_workflow_task_create_trello_card using the `id`. For example:

```terraform
import {
  to = rootly_workflow_task_create_trello_card.my-resource
  id = "00000000-0000-0000-0000-000000000000"
}
```

Locate the resource id in the web app, or retrieve it by listing resources through the API if it's not visible in the web app.

You can generate HCL from the import block using the `-generate-config-out` flag:

```console
terraform plan -generate-config-out=generated.tf
```

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

- `board` (Map of String) Map must contain two fields, `id` and `name`. The board id and display name
- `list` (Map of String) Map must contain two fields, `id` and `name`. The list id and display name
- `title` (String) The card title

Optional:

- `archivation` (Map of String) Map must contain two fields, `id` and `name`. The archivation id and display name
- `description` (String) The card description
- `due_date` (String) The due date
- `labels` (Block List) (see [below for nested schema](#nestedblock--task_params--labels))
- `task_type` (String)

<a id="nestedblock--task_params--labels"></a>
### Nested Schema for `task_params.labels`

Required:

- `id` (String)
- `name` (String)
