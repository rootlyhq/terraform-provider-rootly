resource "rootly_incident_sub_status" "investigating" {
  incident_id   = data.rootly_incident.outage.id
  sub_status_id = rootly_sub_status.investigating.id
  assigned_at   = "2026-06-01T00:00:00Z"
}
