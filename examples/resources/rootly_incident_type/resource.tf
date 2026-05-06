resource "rootly_incident_type" "security_breach" {
  name        = "Security Breach"
  description = "Unauthorized access or data exposure incidents"
  color       = "#FF0000"
}

resource "rootly_incident_type" "service_degradation" {
  name        = "Service Degradation"
  description = "Partial or full service outage affecting customers"
  color       = "#FFA500"
}
