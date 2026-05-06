resource "rootly_on_call_role" "custom" {
  name = "Custom On-Call Role"

  alerts_permissions = [
    "create",
    "update",
    "read",
  ]
  escalation_policies_permissions = [
    "create",
    "read",
    "update",
    "delete",
  ]
  schedules_permissions = [
    "read",
    "update",
  ]
  schedule_override_permissions = [
    "create",
    "read",
    "update",
  ]
}
