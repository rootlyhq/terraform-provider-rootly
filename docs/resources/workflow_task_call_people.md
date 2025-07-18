---
page_title: "Resource rootly_workflow_task_call_people - terraform-provider-rootly"
subcategory: Workflow Tasks
description: |-
    Manages workflow call_people task.
---

# Resource (rootly_workflow_task_call_people)

Manages workflow call_people task.

## Example Usage

```shell
resource "rootly_workflow_incident" "call_people" {
  name        = "Call a teammate"
  description = "Automatically call a specific teammate."
  trigger_params {
    triggers                  = ["incident_created"]
    incident_statuses         = ["started"]
    incident_condition_status = "IS"
  }
  enabled = true
}

resource "rootly_workflow_task_call_people" "call_people" {
  workflow_id     = rootly_workflow_incident.call_people.id
  skip_on_failure = false
  enabled         = true
  task_params {
    name    = "Call people"
    content = "We have an ongoing incident {{ incident.title }} of severity {{ incident.severity }} and your assistance is required."
  }
  # TODO
}
```

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import rootly_workflow_task_call_people using the `id`. For example:

```terraform
import {
  to = rootly_workflow_task_call_people.my-resource
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

- `content` (String) The message to be read by text-to-voice
- `name` (String) The name
- `phone_numbers` (List of String)

Optional:

- `task_type` (String)
