---
page_title: "Resource rootly_workflow_task_invite_to_slack_channel_victor_ops - terraform-provider-rootly"
subcategory: Workflow Tasks
description: |-
    Manages workflow invite_to_slack_channel_victor_ops task.
---

# Resource (rootly_workflow_task_invite_to_slack_channel_victor_ops)

Manages workflow invite_to_slack_channel_victor_ops task.



## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import rootly_workflow_task_invite_to_slack_channel_victor_ops using the `id`. For example:

```terraform
import {
  to = rootly_workflow_task_invite_to_slack_channel_victor_ops.my-resource
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

- `channels` (Block List, Min: 1) (see [below for nested schema](#nestedblock--task_params--channels))
- `team` (Map of String) Map must contain two fields, `id` and `name`.

Optional:

- `task_type` (String)

<a id="nestedblock--task_params--channels"></a>
### Nested Schema for `task_params.channels`

Required:

- `id` (String)
- `name` (String)
