---
page_title: "Resource rootly_retrospective_process - terraform-provider-rootly"
subcategory:
description: |-
    
---

# Resource (rootly_retrospective_process)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the retrospective process
- `retrospective_process_matching_criteria` (Block List, Min: 1, Max: 1) (see [below for nested schema](#nestedblock--retrospective_process_matching_criteria))

### Optional

- `copy_from` (String) Retrospective process ID from which retrospective steps have to be copied. To use starter template for retrospective steps provide value: 'starter_template'
- `description` (String) The description of the retrospective process
- `is_default` (Boolean) Is the retrospective process default?. Value must be one of true or false

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--retrospective_process_matching_criteria"></a>
### Nested Schema for `retrospective_process_matching_criteria`

Optional:

- `group_ids` (List of String) Teams for process matching criteria.
- `incident_type_ids` (List of String) Incident types for process matching criteria.
- `severity_ids` (List of String) Severities for process matching criteria.