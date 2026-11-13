resource "rootly_incident_permission_set" "responders" {
  name                         = "Responders"
  description                  = "Can work private incidents"
  private_incident_permissions = ["read", "update"]
}

# Allow updating the summary of private incidents, but only for SEV0 incidents
# and at incident creation time. severity_params only takes effect when
# severity-based incident permissions are enabled for the team.
resource "rootly_incident_permission_set_boolean" "update_summary" {
  incident_permission_set_id = rootly_incident_permission_set.responders.id
  kind                       = "update_summary"
  private                    = true
  enabled                    = true

  severity_params {
    fully_enabled         = false
    applies_to_unassigned = false
    create_enabled        = true
    severity_ids          = [rootly_severity.sev0.id]
  }
}
