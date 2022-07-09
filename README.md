# Rootly Provider

The [Rootly](https://rootly.com/) provider is used to interact with the resources supported by Rootly. The provider needs to be configured with the proper credentials before it can be used. It requires terraform 0.14 or later.

## Schema

### Optional

- `api_host` (String) The Rootly API host. Defaults to https://api.rootly.com
- `api_token` (String, Sensitive) The Rootly API Token. Generate it from your account at https://rootly.com/account

## Example Usage

```terraform
# Terraform 0.12+ uses the Terraform Registry:

terraform {
  required_providers {
    rootly = {
      source = "rootlyhq/rootly"
    }
  }
}

# Configure the Rootly provider
provider "rootly" {
  # We recommend using the `ROOTLY_API_TOKEN` env var to set the API Token
  # when interacting with Rootly's API.
  # api_token = var.rootly_api_key
}

# Terraform 0.12- can be specified as:
# provider "rootly" {
# We recommend using the `ROOTLY_API_TOKEN` env var to set the API Token
# when interacting with Rootly's API.
# api_token = var.rootly_api_key
# }

resource "rootly_severity" "sev0" {
  name  = "SEV0"
  color = "#FF0000"
}

resource "rootly_service" "elasticsearch_prod" {
  name  = "elasticsearch-prod"
  color = "#800080"
}

resource "rootly_functionality" "add_items_to_cart" {
  name  = "Add items to cart"
  color = "#800080"
}

resource "rootly_workflow_incident" "jira" {
  name = "Create JIRA issue"
  trigger_params {
    triggers = ["incident_created"]
  }
}

resource "rootly_workflow_task_create_jira_issue" "jira" {
  workflow_id = rootly_workflow_incident.jira.id
  task_params {
    title = "{{incident.title}}"
    project_key = "INCIDENTS"
    issue_type = {
      id = "39000"
      name = "Bug"
    }
  }
}
```

## Development

### Auto-generate workflow task resources and tests from Swagger

Use `node tools/gen_tasks.js swagger.json`

The latest Swagger definition can be downloaded from rootly.com:

```sh
curl https://rootly.com/swagger/v1/swagger.json -o ./swagger.json
```

**Note**: The tests for tasks that depend on other resources will fail (ex: rootly_workflow_task_update_incident requires an incident ID). Not all tasks can be reliably auto-generated from Swagger.
