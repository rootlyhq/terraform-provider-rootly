---
page_title: "Resource rootly_escalation_level - terraform-provider-rootly"
subcategory:
description: |-
    
---

# Resource (rootly_escalation_level)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `notification_target_params` (Block List, Min: 1) Escalation level's notification targets (see [below for nested schema](#nestedblock--notification_target_params))
- `position` (Number) Position of the escalation policy level

### Optional

- `delay` (Number) Delay before notification targets will be alerted.
- `escalation_policy_id` (String) The ID of the escalation policy

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--notification_target_params"></a>
### Nested Schema for `notification_target_params`

Required:

- `id` (String)
- `type` (String)