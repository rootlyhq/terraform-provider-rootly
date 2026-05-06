resource "rootly_heartbeat" "nightly_backup" {
  name                     = "nightly-db-backup"
  description              = "Nightly database backup cron job"
  interval                 = 24
  interval_unit            = "hours"
  alert_summary            = "Nightly backup missed"
  notification_target_type = "Service"
  notification_target_id   = rootly_service.elasticsearch_prod.id
  enabled                  = true
}
