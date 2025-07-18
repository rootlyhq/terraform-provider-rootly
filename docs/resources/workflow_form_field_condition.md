---
page_title: "Resource rootly_workflow_form_field_condition - terraform-provider-rootly"
subcategory: Workflows
description: |-
    
---

# Resource (rootly_workflow_form_field_condition)





## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import rootly_workflow_form_field_condition using the `id`. For example:

```terraform
import {
  to = rootly_workflow_form_field_condition.my-resource
  id = "00000000-0000-0000-0000-000000000000"
}
```

Locate the resource id in the web app, or retrieve it by listing resources through the API if it's not visible in the web app.

You can generate HCL from the import block using the `-generate-config-out` flag:

```console
terraform plan -generate-config-out=generated.tf
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `form_field_id` (String) The custom field for this condition

### Optional

- `incident_condition` (String) The trigger condition. Value must be one of `IS`, `ANY`, `CONTAINS`, `CONTAINS_ALL`, `CONTAINS_NONE`, `NONE`, `SET`, `UNSET`.
- `selected_catalog_entity_ids` (List of String)
- `selected_functionality_ids` (List of String)
- `selected_group_ids` (List of String)
- `selected_option_ids` (List of String)
- `selected_service_ids` (List of String)
- `selected_user_ids` (List of Number)
- `values` (List of String)
- `workflow_id` (String) The workflow for this condition

### Read-Only

- `id` (String) The ID of this resource.
