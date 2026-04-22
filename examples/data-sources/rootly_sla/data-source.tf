# Look up an SLA by name
data "rootly_sla" "critical" {
  name = "Critical Incident SLA"
}

# Look up an SLA by slug
data "rootly_sla" "standard" {
  slug = "standard-follow-up-sla"
}

# Use the SLA ID in other resources
output "critical_sla_id" {
  value = data.rootly_sla.critical.id
}
