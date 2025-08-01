---
page_title: "Resource rootly_severity - terraform-provider-rootly"
subcategory:
description: |-
    
---

# Resource (rootly_severity)



## Example Usage

```shell
resource "rootly_severity" "sev0" {
  name          = "SEV0"
  color         = "#FF0000"
  notify_emails = ["foo@acme.com", "bar@acme.com"]
  slack_aliases {
    id   = "S0614TZR7"
    name = "Alias 1" // Any string really
  }
  slack_channels {
    id   = "C06A4RZR9"
    name = "Channel 1" // Any string really
  }
  slack_channels {
    id   = "C02T4RYR2"
    name = "Channel 2" // Any string really
  }
}

resource "rootly_severity" "sev1" {
  name          = "SEV1"
  color         = "#FFA500"
  notify_emails = ["foo@acme.com", "bar@acme.com"]
  slack_aliases {
    id   = "S0614TZR7"
    name = "Alias 1" // Any string really
  }
  slack_channels {
    id   = "C06A4RZR9"
    name = "Channel 1" // Any string really
  }
  slack_channels {
    id   = "C02T4RYR2"
    name = "Channel 2" // Any string really
  }
}

resource "rootly_severity" "sev2" {
  name          = "SEV2"
  color         = "#FFA500"
  notify_emails = ["foo@acme.com", "bar@acme.com"]
  slack_aliases {
    id   = "S0614TZR7"
    name = "Alias 1" // Any string really
  }
  slack_channels {
    id   = "C06A4RZR9"
    name = "Channel 1" // Any string really
  }
  slack_channels {
    id   = "C02T4RYR2"
    name = "Channel 2" // Any string really
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the severity

### Optional

- `color` (String) The hex color of the severity
- `description` (String) The description of the severity
- `notify_emails` (List of String) Emails to attach to the severity
- `position` (Number) Position of the severity
- `severity` (String) The severity of the severity. Value must be one of `critical`, `high`, `medium`, `low`.
- `slack_aliases` (Block List) Slack Aliases associated with this severity (see [below for nested schema](#nestedblock--slack_aliases))
- `slack_channels` (Block List) Slack Channels associated with this severity (see [below for nested schema](#nestedblock--slack_channels))
- `slug` (String) The slug of the severity

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--slack_aliases"></a>
### Nested Schema for `slack_aliases`

Optional:

- `id` (String) Slack alias ID
- `name` (String) Slack alias name


<a id="nestedblock--slack_channels"></a>
### Nested Schema for `slack_channels`

Optional:

- `id` (String) Slack channel ID
- `name` (String) Slack channel name

## Import

rootly_severity can be imported using the [`import` command](https://developer.hashicorp.com/terraform/cli/commands/import).

```sh
terraform import rootly_severity.primary a816421c-6ceb-481a-87c4-585e47451f24
```

Or using an [`import` block](https://developer.hashicorp.com/terraform/language/import).

```terraform
import {
  to = rootly_severity.primary
  id = "a816421c-6ceb-481a-87c4-585e47451f24"
}
```

Locate the resource id in the web app, or retrieve it by listing resources through the API if it's not visible in the web app.

HCL can be generated from the import block using the `-generate-config-out` flag.

```sh
terraform plan -generate-config-out=generated.tf
```
