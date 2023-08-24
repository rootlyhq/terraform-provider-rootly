---
page_title: "Data Source rootly_incident - terraform-provider-rootly"
subcategory:
description: |-
    
---

# Data Source (rootly_incident)



## Example Usage

```terraform
data "rootly_incident" "my-incident" {
  slug = "my-incident-slug"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `acknowledged_at` (Map of String) Filter by date range using 'lt' and 'gt'.
- `created_at` (Map of String) Filter by date range using 'lt' and 'gt'.
- `detected_at` (Map of String) Filter by date range using 'lt' and 'gt'.
- `environments` (String)
- `functionalities` (String)
- `in_triage_at` (Map of String) Filter by date range using 'lt' and 'gt'.
- `kind` (String)
- `mitigated_at` (Map of String) Filter by date range using 'lt' and 'gt'.
- `resolved_at` (Map of String) Filter by date range using 'lt' and 'gt'.
- `services` (String)
- `severity` (String)
- `started_at` (Map of String) Filter by date range using 'lt' and 'gt'.
- `updated_at` (Map of String) Filter by date range using 'lt' and 'gt'.

### Read-Only

- `id` (String) The ID of this resource.