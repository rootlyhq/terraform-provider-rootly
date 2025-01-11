resource "rootly_alert_group" "example" {
  name                 = "Alert group"
  condition_type       = "all"
  time_window          = 10
  group_by_alert_title = true

  attributes {
    json_path = "$.title"
  }

  targets {
    target_type = "Service"
    target_id   = "<Service UUID>"
  }
}
