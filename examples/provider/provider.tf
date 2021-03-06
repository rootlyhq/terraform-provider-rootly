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

# Severities
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

# Services
resource "rootly_service" "elasticsearch_prod" {
  name  = "elasticsearch-prod"
  color = "#800080"
}

resource "rootly_service" "customer_postgresql_prod" {
  name  = "customer-postgresql-prod"
  color = "#800080"
}

# Functionalities
resource "rootly_functionality" "add_items_to_cart" {
  name  = "Add items to cart"
  color = "#800080"
}

resource "rootly_functionality" "logging_in" {
  name  = "Logging In"
  color = "#800080"
}

# Custom Fields
resource "rootly_custom_field" "regions_affected" {
  name = "Regions affected"
  kind = "multi_select"
  shown = ["incident_form"]
  required = ["incident_form"]
}

resource "rootly_custom_field_option" "asia" {
  custom_field_id = rootly_custom_field.regions_affected.id
  value = "Asia"
}

resource "rootly_custom_field_option" "europe" {
  custom_field_id = rootly_custom_field.regions_affected.id
  value = "Europe"
}

resource "rootly_custom_field_option" "north_america" {
  custom_field_id = rootly_custom_field.regions_affected.id
  value = "North America"
}

# Jira workflow
resource "rootly_workflow_incident" "jira" {
  name        = "Create a Jira Issue"
  description = "Open Jira ticket whenever incident starts"
  trigger_params {
    triggers                  = ["incident_created"]
    incident_condition_kind   = "IS"
    incident_kinds            = ["normal"]
    incident_condition_status = "IS"
    incident_statuses         = ["started"]
  }
  enabled = true
}

resource "rootly_workflow_task_create_jira_issue" "jira" {
  workflow_id = rootly_workflow_incident.jira.id
  task_params {
    title       = "{{ incident.title }}"
    description = "{{ incident.summary }}"
    project_key = "ROOT"
    issue_type = {
      id   = "10001"
      name = "Task"
    }
    status = {
      id   = "10000"
      name = "To Do"
    }
    labels = "{{ incident.environment_slugs | concat: incident.service_slugs | concat: incident.functionality_slugs | concat: incident.group_slugs | join: \",\" }}"
  }
}
