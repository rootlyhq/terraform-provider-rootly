resource "rootly_alert_urgency" "critical" {
  name        = "Critical"
  description = "Requires immediate attention — pages on-call"
}

resource "rootly_alert_urgency" "warning" {
  name        = "Warning"
  description = "Should be investigated within 30 minutes"
}

resource "rootly_alert_urgency" "informational" {
  name        = "Informational"
  description = "No immediate action required"
}
