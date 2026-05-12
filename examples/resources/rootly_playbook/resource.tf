resource "rootly_playbook" "database_outage" {
  title        = "Database Outage Response"
  summary      = "Steps to follow when a database outage occurs"
  severity_ids = [rootly_severity.sev0.id, rootly_severity.sev1.id]
  service_ids  = [rootly_service.elasticsearch_prod.id]
}
