---
page_title: "Resource rootly_retrospective_process_group_step - terraform-provider-rootly"
subcategory:
description: |-
    
---

# Resource (rootly_retrospective_process_group_step)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `retrospective_step_id` (String)

### Optional

- `position` (Number)
- `retrospective_process_group_id` (String)

### Read-Only

- `id` (String) The ID of this resource.

## Import

rootly_retrospective_process_group_step can be imported using the [`import` command](https://developer.hashicorp.com/terraform/cli/commands/import).

```sh
terraform import rootly_retrospective_process_group_step.primary a816421c-6ceb-481a-87c4-585e47451f24
```

Or using an [`import` block](https://developer.hashicorp.com/terraform/language/import).

```terraform
import {
  to = rootly_retrospective_process_group_step.primary
  id = "a816421c-6ceb-481a-87c4-585e47451f24"
}
```

Locate the resource id in the web app, or retrieve it by listing resources through the API if it's not visible in the web app.

HCL can be generated from the import block using the `-generate-config-out` flag.

```sh
terraform plan -generate-config-out=generated.tf
```
