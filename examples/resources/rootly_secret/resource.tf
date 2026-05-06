resource "rootly_secret" "pagerduty_api_token" {
  name   = "pagerduty_api_token"
  secret = var.pagerduty_api_token
}
