---
page_title: "Resource rootly_incident_permission_set_resource - terraform-provider-rootly"
subcategory:
description: |-
    
---

# Resource (rootly_incident_permission_set_resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `incident_permission_set_id` (String)

### Optional

- `kind` (String) Value must be one of `severities`, `incident_types`, `statuses`, `sub_statuses`.
- `private` (Boolean) Value must be one of true or false
- `resource_id` (String)
- `resource_type` (String)

### Read-Only

- `id` (String) The ID of this resource.
