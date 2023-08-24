---
page_title: "Resource rootly_workflow_task_update_notion_page - terraform-provider-rootly"
subcategory: Workflow Tasks
description: |-
    Manages workflow updatenotionpage task.
---

# Resource (rootly_workflow_task_update_notion_page)

Manages workflow update_notion_page task.



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

- `file_id` (String) The Notion page ID

Optional:

- `post_mortem_template_id` (String) Retrospective template to use when creating page task, if desired.
- `show_action_items_as_table` (Boolean)
- `show_timeline_as_table` (Boolean)
- `task_type` (String)