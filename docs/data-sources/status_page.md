---
page_title: "Data Source rootly_status_page - terraform-provider-rootly"
subcategory:
description: |-
    
---

# Data Source (rootly_status_page)



## Example Usage

```shell
data "rootly_status_page" "my-status-page" {
  slug = "my-status-page"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `created_at` (Map of String) Filter by date range using 'lt' and 'gt'.
- `slug` (String)

### Read-Only

- `id` (String) The ID of this resource.
