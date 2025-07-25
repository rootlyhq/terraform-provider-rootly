---
page_title: "Resource rootly_custom_form - terraform-provider-rootly"
subcategory:
description: |-
    
---

# Resource (rootly_custom_form)





## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import rootly_custom_form using the `id`. For example:

```terraform
import {
  to = rootly_custom_form.my-resource
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

- `command` (String) The Slack command used to trigger this form.
- `name` (String) The name of the custom form.

### Optional

- `description` (String)
- `enabled` (Boolean)
- `slug` (String) The custom form slug. Add this to form_field.shown or form_field.required to associate form fields with custom forms.

### Read-Only

- `id` (String) The ID of this resource.
