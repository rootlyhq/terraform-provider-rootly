---
page_title: "Resource rootly_workflow_task_page_pagerduty_on_call_responders - terraform-provider-rootly"
subcategory: ""
description: |-
    Manages workflow pagepagerdutyoncallresponders task.
---

# Resource (rootly_workflow_task_page_pagerduty_on_call_responders)

Manages workflow page_pagerduty_on_call_responders task.

## Example Usage

```
resource "rootly_workflow_task_page_pagerduty_on_call_responders" "foo" {
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

- `service` (Map of String)

Optional:

- `escalation_policies` (Block List) (see [below for nested schema](#nestedblock--task_params--escalation_policies))
- `message` (String)
- `task_type` (String)
- `urgency` (String)
- `users` (Block List) (see [below for nested schema](#nestedblock--task_params--users))

<a id="nestedblock--task_params--escalation_policies"></a>
### Nested Schema for `task_params.escalation_policies`

Required:

- `name` (String)

Read-Only:

- `id` (String) The ID of this resource.


<a id="nestedblock--task_params--users"></a>
### Nested Schema for `task_params.users`

Required:

- `name` (String)

Read-Only:

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import rootly_workflow_task_page_pagerduty_on_call_responders.foo 11111111-2222-3333-4444-555555555555
```