resource "rootly_incident_sub_status" "investigating" {
  incident_id   = rootly_incident.outage.id
  sub_status_id = rootly_sub_status.investigating.id
}
