---
page_title: "Resource rootly_custom_field - terraform-provider-rootly"
subcategory: ""
description: |-
    Manages custom fields.
---

# Resource (rootly_custom_field)

Manages custom fields.

## Example Usage

```
resource "rootly_custom_field" "foo" {
  name = "bar"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `kind` (String) The kind of the custom field
- `label` (String) The name of the custom field

### Optional

- `description` (String) The description of the custom field
- `enabled` (Boolean) Whether the custom field is enabled or not
- `required` (List of String) Where the custom field is required.
- `shown` (List of String) Where the custom field is shown.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import rootly_custom_field.foo 11111111-2222-3333-4444-555555555555
```