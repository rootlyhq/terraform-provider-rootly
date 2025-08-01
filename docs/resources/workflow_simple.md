---
page_title: "Resource rootly_workflow_simple - terraform-provider-rootly"
subcategory: Workflows
description: |-
    
---

# Resource (rootly_workflow_simple)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The title of the workflow

### Optional

- `cause_ids` (List of String)
- `command` (String) Workflow command
- `command_feedback_enabled` (Boolean) This will notify you back when the workflow is starting. Value must be one of true or false
- `continuously_repeat` (Boolean) When continuously repeat is true, repeat workflows aren't automatically stopped when conditions aren't met. This setting won't override your conditions set by repeat_condition_duration_since_first_run and repeat_condition_number_of_repeats parameters.. Value must be one of true or false
- `description` (String) The description of the workflow
- `enabled` (Boolean)
- `environment_ids` (List of String)
- `functionality_ids` (List of String)
- `group_ids` (List of String)
- `incident_role_ids` (List of String)
- `incident_type_ids` (List of String)
- `locked` (Boolean) Restricts workflow edits to admins when turned on. Only admins can set this field.. Value must be one of true or false
- `position` (Number) The order which the workflow should run with other workflows.
- `repeat_condition_duration_since_first_run` (String) The workflow will stop repeating if its runtime since it's first workflow run exceeds the duration set in this field
- `repeat_condition_number_of_repeats` (Number) The workflow will stop repeating if the number of repeats exceeds the value set in this field
- `repeat_every_duration` (String) Repeat workflow every duration
- `repeat_on` (List of String) Repeat on weekdays. Value must be one of `S`, `M`, `T`, `W`, `R`, `F`, `U`.
- `service_ids` (List of String)
- `severity_ids` (List of String)
- `slug` (String) The slug of the workflow
- `sub_status_ids` (List of String)
- `trigger_params` (Block List, Max: 1) (see [below for nested schema](#nestedblock--trigger_params))
- `wait` (String) Wait this duration before executing
- `workflow_group_id` (String) The group this workflow belongs to.

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--trigger_params"></a>
### Nested Schema for `trigger_params`

Optional:

- `trigger_type` (String) Value must be one off `simple`.
- `triggers` (List of String) Actions that trigger the workflow. Value must be one of `slack_command`.

## Import

rootly_workflow_simple can be imported using the [`import` command](https://developer.hashicorp.com/terraform/cli/commands/import).

```sh
terraform import rootly_workflow_simple.primary a816421c-6ceb-481a-87c4-585e47451f24
```

Or using an [`import` block](https://developer.hashicorp.com/terraform/language/import).

```terraform
import {
  to = rootly_workflow_simple.primary
  id = "a816421c-6ceb-481a-87c4-585e47451f24"
}
```

Locate the resource id in the web app, or retrieve it by listing resources through the API if it's not visible in the web app.

HCL can be generated from the import block using the `-generate-config-out` flag.

```sh
terraform plan -generate-config-out=generated.tf
```
