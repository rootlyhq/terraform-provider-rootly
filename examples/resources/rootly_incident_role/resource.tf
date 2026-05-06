resource "rootly_incident_role" "incident_commander" {
  name        = "Incident Commander"
  summary     = "Leads the incident response"
  description = "Coordinates response, delegates tasks, makes decisions"
  enabled     = true
}

resource "rootly_incident_role" "communications_lead" {
  name     = "Communications Lead"
  summary  = "Manages stakeholder communications"
  optional = true
  enabled  = true
}
