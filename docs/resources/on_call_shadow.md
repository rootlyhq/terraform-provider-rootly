---
page_title: "Resource rootly_on_call_shadow - terraform-provider-rootly"
subcategory:
description: |-
    
---

# Resource (rootly_on_call_shadow)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `ends_at` (String) End datetime for shadow shift
- `shadow_user_id` (Number) Which user the shadow shift belongs to.
- `shadowable_id` (String) ID of schedule or user the shadow user is shadowing
- `starts_at` (String) Start datetime of shadow shift

### Optional

- `schedule_id` (String) ID of schedule the shadow shift belongs to
- `shadowable_type` (String) Value must be one of `User`, `Schedule`.

### Read-Only

- `id` (String) The ID of this resource.