---
page_title: "Resource rootly_workflow_task_update_incident - terraform-provider-rootly"
subcategory: Workflow Tasks
description: |-
    Manages workflow update_incident task.
---

# Resource (rootly_workflow_task_update_incident)

Manages workflow update_incident task.



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

- `incident_id` (String) The incident id to update or id of any attribute on the incident

Optional:

- `acknowledged_at` (String)
- `attribute_to_query_by` (String) Value must be one of `id`, `slug`, `sequential_id`, `pagerduty_incident_id`, `opsgenie_incident_id`, `victor_ops_incident_id`, `jira_issue_id`, `asana_task_id`, `shortcut_task_id`, `linear_issue_id`, `zendesk_ticket_id`, `motion_task_id`, `trello_card_id`, `airtable_record_id`, `shortcut_story_id`, `github_issue_id`, `gitlab_issue_id`, `freshservice_ticket_id`, `freshservice_task_id`, `clickup_task_id`.
- `custom_fields_mapping` (String) Custom field mappings. Can contain liquid markup and need to be valid JSON
- `detected_at` (String)
- `environment_ids` (List of String)
- `functionality_ids` (List of String)
- `group_ids` (List of String)
- `incident_type_ids` (List of String)
- `mitigated_at` (String)
- `private` (Boolean) Value must be one of true or false
- `resolved_at` (String)
- `service_ids` (List of String)
- `severity_id` (String)
- `started_at` (String)
- `status` (String)
- `summary` (String) The incident summary
- `task_type` (String)
- `title` (String) The incident title

## Import

rootly_workflow_task_update_incident can be imported using the [`import` command](https://developer.hashicorp.com/terraform/cli/commands/import).

```sh
terraform import rootly_workflow_task_update_incident.primary a816421c-6ceb-481a-87c4-585e47451f24
```

Or using an [`import` block](https://developer.hashicorp.com/terraform/language/import).

```terraform
import {
  to = rootly_workflow_task_update_incident.primary
  id = "a816421c-6ceb-481a-87c4-585e47451f24"
}
```

Locate the resource id in the web app, or retrieve it by listing resources through the API if it's not visible in the web app.

HCL can be generated from the import block using the `-generate-config-out` flag.

```sh
terraform plan -generate-config-out=generated.tf
```
