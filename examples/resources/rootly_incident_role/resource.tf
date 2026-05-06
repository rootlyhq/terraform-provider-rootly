resource "rootly_incident_role" "commander" {
  name        = "Incident Commander"
  summary     = "Leads the incident response"
  description = "Responsible for coordinating the incident response effort"
  enabled     = true
}

resource "rootly_incident_role" "communications_lead" {
  name        = "Communications Lead"
  summary     = "Manages stakeholder communications"
  description = "Responsible for internal and external communications during the incident"
  optional    = true
  enabled     = true
}
