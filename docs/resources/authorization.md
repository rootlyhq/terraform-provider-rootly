---
page_title: "Resource rootly_authorization - terraform-provider-rootly"
subcategory:
description: |-
    
---

# Resource (rootly_authorization)





## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import rootly_authorization using the `id`. For example:

```terraform
import {
  to = rootly_authorization.my-resource
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

- `authorizable_id` (String) The id of the resource being accessed.
- `grantee_id` (String) The resource id granted access.
- `permissions` (List of String) Value must be one of `read`, `update`, `authorize`, `destroy`.

### Optional

- `authorizable_type` (String) The type of resource being accessed.. Value must be one of `Dashboard`.
- `grantee_type` (String) The type of resource granted access.. Value must be one of `User`, `Team`.

### Read-Only

- `id` (String) The ID of this resource.
