resource "rootly_webhooks_endpoint" "datadog" {
  name        = "Datadog Webhook"
  url         = "https://app.datadoghq.com/api/v1/events"
  enabled     = true
  event_types = ["incident.created", "incident.resolved"]
}
