---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "rootly_custom_fields Data Source - terraform-provider-rootly"
subcategory: ""
description: |-
  
---

# rootly_custom_fields (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `enabled` (Boolean)
- `kind` (String)
- `label` (String)
- `slug` (String)

### Read-Only

- `custom_fields` (List of Object) (see [below for nested schema](#nestedatt--custom_fields))
- `id` (String) The ID of this resource.

<a id="nestedatt--custom_fields"></a>
### Nested Schema for `custom_fields`

Read-Only:

- `description` (String)
- `enabled` (Boolean)
- `id` (String)
- `kind` (String)
- `label` (String)
- `required` (List of String)
- `shown` (List of String)
- `slug` (String)

