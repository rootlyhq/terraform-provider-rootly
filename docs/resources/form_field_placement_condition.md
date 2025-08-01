---
page_title: "Resource rootly_form_field_placement_condition - terraform-provider-rootly"
subcategory:
description: |-
    
---

# Resource (rootly_form_field_placement_condition)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `form_field_id` (String) The condition field.
- `values` (List of String) The values for comparison.

### Optional

- `comparison` (String) The condition comparison.. Value must be one of `equal`, `is_set`, `is_not_set`.
- `conditioned` (String) The resource or attribute the condition applies.. Value must be one of `placement`, `required`.
- `form_field_placement_id` (String) The form field placement this condition applies.
- `position` (Number) The condition position.

### Read-Only

- `id` (String) The ID of this resource.

## Import

rootly_form_field_placement_condition can be imported using the [`import` command](https://developer.hashicorp.com/terraform/cli/commands/import).

```sh
terraform import rootly_form_field_placement_condition.primary a816421c-6ceb-481a-87c4-585e47451f24
```

Or using an [`import` block](https://developer.hashicorp.com/terraform/language/import).

```terraform
import {
  to = rootly_form_field_placement_condition.primary
  id = "a816421c-6ceb-481a-87c4-585e47451f24"
}
```

Locate the resource id in the web app, or retrieve it by listing resources through the API if it's not visible in the web app.

HCL can be generated from the import block using the `-generate-config-out` flag.

```sh
terraform plan -generate-config-out=generated.tf
```
