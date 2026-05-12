resource "rootly_playbook_task" "check_dashboards" {
  playbook_id = rootly_playbook.database_outage.id
  task        = "Check Grafana dashboards for anomalies"
}

resource "rootly_playbook_task" "notify_dba" {
  playbook_id = rootly_playbook.database_outage.id
  task        = "Notify DBA on-call via Slack"
}
