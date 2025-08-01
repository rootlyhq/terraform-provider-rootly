---
page_title: "Resource rootly_status_page - terraform-provider-rootly"
subcategory:
description: |-
    
---

# Resource (rootly_status_page)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `title` (String) The title of the status page

### Optional

- `allow_search_engine_index` (Boolean) Allow search engines to include your public status page in search results. Value must be one of true or false
- `authentication_enabled` (Boolean) Enable authentication. Value must be one of true or false
- `authentication_password` (String) Authentication password
- `description` (String) The description of the status page
- `enabled` (Boolean)
- `external_domain_names` (List of String) External domain names attached to the status page
- `failure_message` (String) Message showing when at least one component is not operational
- `footer_color` (String) The color of the footer. Eg. "#1F2F41"
- `functionality_ids` (List of String) Functionalities attached to the status page
- `ga_tracking_id` (String) Google Analytics tracking ID
- `header_color` (String) The color of the header. Eg. "#0061F2"
- `public` (Boolean) Make the status page accessible to the public. Value must be one of true or false
- `public_description` (String) The public description of the status page
- `public_title` (String) The public title of the status page
- `service_ids` (List of String) Services attached to the status page
- `show_uptime` (Boolean) Show uptime. Value must be one of true or false
- `show_uptime_last_days` (Number) Show uptime over x days. Value must be one of `30`, `60`, `90`.
- `slug` (String) The slug of the status page
- `success_message` (String) Message showing when all components are operational
- `time_zone` (String) A valid IANA time zone name.
- `website_privacy_url` (String) Website Privacy URL
- `website_support_url` (String) Website Support URL
- `website_url` (String) Website URL

### Read-Only

- `id` (String) The ID of this resource.

## Import

rootly_status_page can be imported using the [`import` command](https://developer.hashicorp.com/terraform/cli/commands/import).

```sh
terraform import rootly_status_page.primary a816421c-6ceb-481a-87c4-585e47451f24
```

Or using an [`import` block](https://developer.hashicorp.com/terraform/language/import).

```terraform
import {
  to = rootly_status_page.primary
  id = "a816421c-6ceb-481a-87c4-585e47451f24"
}
```

Locate the resource id in the web app, or retrieve it by listing resources through the API if it's not visible in the web app.

HCL can be generated from the import block using the `-generate-config-out` flag.

```sh
terraform plan -generate-config-out=generated.tf
```
