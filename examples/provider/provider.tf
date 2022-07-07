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

resource "rootly_severity" "sev1" {
  name  = "SEV1"
  color = "#FFA500"
}

resource "rootly_severity" "sev2" {
  name  = "SEV2"
  color = "#FFA500"
}

resource "rootly_service" "elasticsearch_prod" {
  name  = "elasticsearch-prod"
  color = "#800080"
}

resource "rootly_service" "customer_postgresql_prod" {
  name  = "customer-postgresql-prod"
  color = "#800080"
}

resource "rootly_functionality" "add_items_to_cart" {
  name  = "Add items to cart"
  color = "#800080"
}

resource "rootly_functionality" "logging_in" {
  name  = "Logging In"
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
