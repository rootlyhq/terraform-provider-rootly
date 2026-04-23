# Look up existing resources to reference in conditions
data "rootly_severity" "sev0" {
  slug = "sev0"
}

data "rootly_severity" "sev1" {
  slug = "sev1"
}

data "rootly_incident_role" "commander" {
  slug = "incident-commander"
}

# Basic SLA — follow-ups must be assigned within 3 days of the incident
# starting, and completed within 7 days of resolution.
resource "rootly_sla" "basic" {
  name                              = "Standard Follow-Up SLA"
  description                       = "Ensure follow-ups are assigned and completed on time"
  assignment_deadline_days          = 3
  assignment_deadline_parent_status = "started"
  completion_deadline_days          = 7
  completion_deadline_parent_status = "resolved"
  manager_role_id                   = data.rootly_incident_role.commander.id
}

# SLA with conditions — only applies to SEV0/SEV1 incidents.
# Note: the `values` field takes resource IDs (not display names like "SEV0").
# Use a data source to look up the ID first.
resource "rootly_sla" "critical" {
  name                              = "Critical Incident SLA"
  description                       = "Stricter deadlines for SEV0 and SEV1 incidents"
  condition_match_type              = "ALL"
  assignment_deadline_days          = 3
  assignment_deadline_parent_status = "started"
  assignment_skip_weekends          = false
  completion_deadline_days          = 5
  completion_deadline_parent_status = "resolved"
  completion_skip_weekends          = false
  manager_role_id                   = data.rootly_incident_role.commander.id

  # is_one_of: match incidents with SEV0 or SEV1 severity.
  # Note: `values` takes resource IDs (not display names like "SEV0").
  # Use a data source to look up the ID first.
  conditions {
    conditionable_type = "SLAs::BuiltInFieldCondition"
    property           = "severity"
    operator           = "is_one_of"
    values             = [data.rootly_severity.sev0.id, data.rootly_severity.sev1.id]
  }

  # is_set: require that the environment field is present (no values needed)
  conditions {
    conditionable_type = "SLAs::BuiltInFieldCondition"
    property           = "environment"
    operator           = "is_set"
  }

  # Notify 1 day before the deadline, when due, and 1 day after
  notification_configurations {
    offset_type = "before_due"
    offset_days = 1
  }

  notification_configurations {
    offset_type = "when_due"
    offset_days = 0
  }

  notification_configurations {
    offset_type = "after_due"
    offset_days = 1
  }
}

# SLA with a custom field condition using the "contains" operator (single value)
resource "rootly_sla" "compliance" {
  name                              = "Compliance Review SLA"
  assignment_deadline_days          = 2
  assignment_deadline_parent_status = "started"
  completion_deadline_days          = 5
  completion_deadline_parent_status = "resolved"
  manager_role_id                   = data.rootly_incident_role.commander.id

  # contains: match when the custom field value contains a substring
  conditions {
    conditionable_type = "SLAs::CustomFieldCondition"
    form_field_id      = "your-custom-field-uuid"
    operator           = "contains"
    values             = ["production"]
  }
}
