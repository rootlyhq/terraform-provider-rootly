resource "rootly_incident_type" "security" {
  name        = "Security"
  description = "Security-related incidents"
  color       = "#FF0000"
}

resource "rootly_incident_type" "infrastructure" {
  name        = "Infrastructure"
  description = "Infrastructure and platform incidents"
  color       = "#FFA500"
}
