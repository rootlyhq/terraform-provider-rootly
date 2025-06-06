---
page_title: "Resource rootly_custom_field_option - terraform-provider-rootly"
subcategory:
description: |-
    DEPRECATED: Please use rootly_form_field and rootly_form_field_option resources instead.
---

# Resource (rootly_custom_field_option)

DEPRECATED: Please use `rootly_form_field` and `rootly_form_field_option` resources instead.



## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import rootly_custom_field_option using the `id`. For example:

```terraform
import {
  to rootly_custom_field_option.my-resource
  id = "00000000-0000-0000-0000-000000000000"
}
```

Using `terraform import`, import rootly_custom_field_option using the `id`. For example:

```console
% terraform import rootly_custom_field_option.my-resource 00000000-0000-0000-0000-000000000000
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `value` (String) The value of the custom_field_option

### Optional

- `color` (String) The hex color of the custom_field_option
- `custom_field_id` (Number) The ID of the parent custom field
- `default` (Boolean)
- `position` (Number) The position of the custom_field_option

### Read-Only

- `id` (String) The ID of this resource.
