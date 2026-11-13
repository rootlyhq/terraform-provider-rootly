resource "rootly_incident_permission_set" "responders" {
  name                        = "Responders"
  description                 = "Can work public incidents"
  public_incident_permissions = ["read", "update"]
}

# Allow setting the SEV1 severity on public incidents, but only on incidents
# that are already SEV0 or SEV1. severity_params only takes effect when
# severity-based incident permissions are enabled for the team.
resource "rootly_incident_permission_set_resource" "sev1" {
  incident_permission_set_id = rootly_incident_permission_set.responders.id
  kind                       = "severities"
  private                    = false
  resource_type              = "Severity"
  resource_id                = rootly_severity.sev1.id

  severity_params {
    fully_enabled         = false
    applies_to_unassigned = false
    severity_ids          = [rootly_severity.sev0.id, rootly_severity.sev1.id]
  }
}
