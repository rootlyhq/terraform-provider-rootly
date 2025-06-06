---
page_title: "Resource rootly_form_field_placement - terraform-provider-rootly"
subcategory:
description: |-
    
---

# Resource (rootly_form_field_placement)





## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import rootly_form_field_placement using the `id`. For example:

```terraform
import {
  to rootly_form_field_placement.my-resource
  id = "00000000-0000-0000-0000-000000000000"
}
```

Using `terraform import`, import rootly_form_field_placement using the `id`. For example:

```console
% terraform import rootly_form_field_placement.my-resource 00000000-0000-0000-0000-000000000000
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `form` (String) The form this field is placed on.
- `form_set_id` (String) The form set this field is placed in.

### Optional

- `form_field_id` (String) The form field that is placed.
- `placement_operator` (String) Logical operator when evaluating multiple form_field_placement_conditions with conditioned=placement. Value must be one of `and`, `or`.
- `position` (Number) The position of the field placement.
- `required` (Boolean) Whether the field is unconditionally required on this form.. Value must be one of true or false
- `required_operator` (String) Logical operator when evaluating multiple form_field_placement_conditions with conditioned=required. Value must be one of `and`, `or`.

### Read-Only

- `id` (String) The ID of this resource.
