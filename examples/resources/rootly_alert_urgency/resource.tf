resource "rootly_alert_urgency" "critical" {
  name        = "Critical"
  description = "Requires immediate attention"
}

resource "rootly_alert_urgency" "warning" {
  name        = "Warning"
  description = "Should be investigated soon"
}

resource "rootly_alert_urgency" "info" {
  name        = "Informational"
  description = "No immediate action required"
}
