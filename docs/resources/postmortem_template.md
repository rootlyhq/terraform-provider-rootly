---
page_title: "Resource rootly_postmortem_template - terraform-provider-rootly"
subcategory:
description: |-
    
---

# Resource (rootly_postmortem_template)



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `content` (String) The postmortem template. Liquid syntax and markdown are supported.
- `name` (String) The name of the postmortem template

### Optional

- `default` (Boolean) Default selected template when editing a postmortem

### Read-Only

- `id` (String) The ID of this resource.