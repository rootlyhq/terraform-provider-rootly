---
page_title: "Resource rootly_dashboard_panel - terraform-provider-rootly"
subcategory:
description: |-
    Manages dashboard_panels.
---

# Resource (rootly_dashboard_panel)

Manages dashboard_panels.

## Example Usage

```shell
resource "rootly_dashboard" "overview" {
  name = "my_dashboard"
}

resource "rootly_dashboard_panel" "incidents_by_severity" {
  dashboard_id = rootly_dashboard.foo.id
  name         = "test"
  params {
    display = "line_chart"
    datasets {
      collection = "incidents"
      filter {
        operation = "and"
        rules {
          operation = "and"
          condition = "="
          key       = "status"
          value     = "started"
        }
      }
      group_by = "severity"
      aggregate {
        cumulative = false
        key        = "results"
        operation  = "count"
      }
    }
  }
}
```

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import rootly_dashboard_panel using the `id`. For example:

```terraform
import {
  to = rootly_dashboard_panel.my-resource
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

- `dashboard_id` (String) The id of the parent dashboard
- `params` (Block List, Min: 1, Max: 1) The params JSON of the dashboard_panel. See rootly API docs for schema. (see [below for nested schema](#nestedblock--params))

### Optional

- `name` (String) The name of the dashboard_panel
- `position` (Block List, Max: 1) (see [below for nested schema](#nestedblock--position))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--params"></a>
### Nested Schema for `params`

Required:

- `display` (String)

Optional:

- `datasets` (Block List) (see [below for nested schema](#nestedblock--params--datasets))
- `legend` (Block List, Max: 1) (see [below for nested schema](#nestedblock--params--legend))

<a id="nestedblock--params--datasets"></a>
### Nested Schema for `params.datasets`

Required:

- `collection` (String)

Optional:

- `aggregate` (Block List, Max: 1) (see [below for nested schema](#nestedblock--params--datasets--aggregate))
- `filter` (Block List) (see [below for nested schema](#nestedblock--params--datasets--filter))
- `group_by` (String)
- `name` (String)

<a id="nestedblock--params--datasets--aggregate"></a>
### Nested Schema for `params.datasets.aggregate`

Required:

- `cumulative` (Boolean)
- `key` (String)
- `operation` (String)


<a id="nestedblock--params--datasets--filter"></a>
### Nested Schema for `params.datasets.filter`

Required:

- `operation` (String)

Optional:

- `rules` (Block List) (see [below for nested schema](#nestedblock--params--datasets--filter--rules))

<a id="nestedblock--params--datasets--filter--rules"></a>
### Nested Schema for `params.datasets.filter.rules`

Required:

- `condition` (String)
- `key` (String)
- `operation` (String)
- `value` (String)




<a id="nestedblock--params--legend"></a>
### Nested Schema for `params.legend`

Required:

- `groups` (String)



<a id="nestedblock--position"></a>
### Nested Schema for `position`

Required:

- `h` (Number)
- `w` (Number)
- `x` (Number)
- `y` (Number)
