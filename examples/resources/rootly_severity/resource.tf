resource "rootly_severity" "sev0" {
  name          = "SEV0"
  color         = "#FF0000"
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

resource "rootly_severity" "sev1" {
  name          = "SEV1"
  color         = "#FFA500"
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

resource "rootly_severity" "sev2" {
  name          = "SEV2"
  color         = "#FFA500"
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
