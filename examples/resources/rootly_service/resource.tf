resource "rootly_service" "elasticsearch_prod" {
  name          = "elasticsearch-prod"
  color         = "#800080"
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

resource "rootly_service" "customer_postgresql_prod" {
  name          = "customer-postgresql-prod"
  color         = "#800080"
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
