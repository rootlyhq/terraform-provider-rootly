resource "rootly_heartbeat" "cron_job" {
  name                     = "nightly-backup"
  description              = "Nightly database backup job"
  interval                 = 24
  interval_unit            = "hours"
  alert_summary            = "Nightly backup missed"
  notification_target_type = "Service"
  notification_target_id   = rootly_service.database.id
  enabled                  = true
}
