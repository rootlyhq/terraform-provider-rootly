---
page_title: "Resource rootly_post_mortem_template - terraform-provider-rootly"
subcategory:
description: |-
    
---

# Resource (rootly_post_mortem_template)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the postmortem template

### Optional

- `content` (String) The postmortem template. Liquid syntax and markdown are supported
- `default` (Boolean) Default selected template when editing a postmortem. Value must be one of true or false
- `format` (String) The format of the input. Value must be one of `html`, `markdown`.

### Read-Only

- `id` (String) The ID of this resource.
