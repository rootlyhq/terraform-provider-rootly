# Terraform 0.12+ uses the Terraform Registry:

terraform {
  required_providers {
    rootly = {
      source = "rootly/rootly"
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
