---
page_title: "Resource rootly_sub_status - terraform-provider-rootly"
subcategory:
description: |-
    
---

# Resource (rootly_sub_status)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String)

### Optional

- `description` (String)
- `parent_status` (String) Value must be one of `in_triage`, `started`, `resolved`, `closed`, `cancelled`, `scheduled`, `in_progress`, `completed`.
- `position` (Number)
- `slug` (String)

### Read-Only

- `id` (String) The ID of this resource.