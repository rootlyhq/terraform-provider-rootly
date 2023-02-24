resource "rootly_functionality" "add_items_to_card" {
  name          = "Add items to cart"
  color         = "#800080"
  notify_emails = ["foo@acme.com", "bar@acme.com"]
  slack_aliases = [{
    id   = "S0614TZR7"
    name = "Alias 1" // Any string really
  }]
  slack_channels = [{
    id   = "C06A4RZR9"
    name = "Channel 1" // Any string really,
    },
    {
      id   = "C02T4RYR2"
      name = "Channel 2" // Any string really,
    }
  ]
}

resource "rootly_functionality" "logging_in" {
  name          = "Logging In"
  color         = "#800080"
  notify_emails = ["foo@acme.com", "bar@acme.com"]
  slack_aliases = [{
    id   = "S0614TZR7"
    name = "Alias 1" // Any string really
  }]
  slack_channels = [{
    id   = "C06A4RZR9"
    name = "Channel 1" // Any string really,
    },
    {
      id   = "C02T4RYR2"
      name = "Channel 2" // Any string really,
    }
  ]
}
