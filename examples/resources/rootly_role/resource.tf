resource "rootly_role" "custom" {
  name = "Custom Rootly Role"

  incidents_permissions = [
    "read",
    "update",
  ]
  incident_feedbacks_permissions = [
    "create",
    "read",
  ]
  services_permissions = [
    "read",
  ]
  severities_permissions = [
    "read",
  ]

  is_editable  = true
  is_deletable = true
}
