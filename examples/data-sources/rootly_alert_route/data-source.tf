# Find an alert route by name
data "rootly_alert_route" "example" {
  name = "Production Alerts"
}

# Use the alert route ID in other resources
output "alert_route_id" {
  value = data.rootly_alert_route.example.id
}