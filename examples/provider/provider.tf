terraform {
  required_providers {
    rootly = {
      source = "rootlyhq/rootly"
    }
  }
}

provider "rootly" {
  # We recommend using the `ROOTLY_API_TOKEN` env var to set the API Token
  # when interacting with Rootly's API.
  # api_token = var.rootly_api_key
}
