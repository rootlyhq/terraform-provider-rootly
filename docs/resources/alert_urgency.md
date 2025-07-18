---
page_title: "Resource rootly_alert_urgency - terraform-provider-rootly"
subcategory:
description: |-
    
---

# Resource (rootly_alert_urgency)





## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import rootly_alert_urgency using the `id`. For example:

```terraform
import {
  to = rootly_alert_urgency.my-resource
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

- `description` (String) The description of the alert urgency
- `name` (String) The name of the alert urgency

### Optional

- `position` (Number) Position of the alert urgency

### Read-Only

- `id` (String) The ID of this resource.
