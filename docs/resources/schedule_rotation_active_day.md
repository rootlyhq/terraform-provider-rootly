---
page_title: "Resource rootly_schedule_rotation_active_day - terraform-provider-rootly"
subcategory:
description: |-
    
---

# Resource (rootly_schedule_rotation_active_day)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `active_time_attributes` (Block List, Min: 1) Schedule rotation active times per day (see [below for nested schema](#nestedblock--active_time_attributes))

### Optional

- `day_name` (String) Schedule rotation day name for which active times to be created. Value must be one of `S`, `M`, `T`, `W`, `R`, `F`, `U`.
- `schedule_rotation_id` (String)

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--active_time_attributes"></a>
### Nested Schema for `active_time_attributes`

Required:

- `id` (String)
- `name` (String)