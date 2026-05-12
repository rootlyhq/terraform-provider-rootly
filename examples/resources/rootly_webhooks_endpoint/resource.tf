resource "rootly_webhooks_endpoint" "datadog_events" {
  name        = "Datadog Events"
  url         = "https://app.datadoghq.com/api/v1/events"
  enabled     = true
  event_types = ["incident.created", "incident.resolved"]
}
