data "rootly_user" "user_1" {
  email = "john@example.com"
}

# Create teams and escalation policies to set as destinations/owner teams
resource "rootly_team" "on_call_team" {
  name     = "On-Call Team"
  user_ids = [data.rootly_user.user_1.id]
}

resource "rootly_team" "security_team" {
  name     = "Security Team"
  user_ids = [data.rootly_user.user_1.id]
}

resource "rootly_escalation_policy" "production_ep" {
  name      = "Production Escalation"
  group_ids = [rootly_team.on_call_team.id]
}

resource "rootly_escalation_policy" "security_ep" {
  name      = "Security Escalation"
  group_ids = [rootly_team.security_team.id]
}

# Create alert fields for rule conditions
resource "rootly_alert_field" "severity_field" {
  name = "Severity"
}

resource "rootly_alert_field" "service_field" {
  name = "Service"
}

# Create alert source to route from
resource "rootly_alerts_source" "monitoring_source" {
  name        = "Production Monitoring"
  source_type = "generic_webhook"

  alert_source_fields_attributes {
    alert_field_id = "538914bc-27d5-40d4-944b-ab1d71cc59b6" # Built-in title field
    template_body  = "{{ alert.title }}"
  }

  alert_source_fields_attributes {
    alert_field_id = resource.rootly_alert_field.severity_field.id
    template_body  = "{{ alert.severity }}"
  }
}

# Create alert route with multiple rules
resource "rootly_alert_route" "production_route" {
  name              = "Production Alert Routing"
  alerts_source_ids = [resource.rootly_alerts_source.monitoring_source.id]
  owning_team_ids   = [resource.rootly_team.on_call_team.id]
  enabled           = true

  rules {
    name = "High Severity Route"
    destinations {
      target_type = "EscalationPolicy"
      target_id   = rootly_escalation_policy.production_ep.id
    }
    condition_groups {
      position = 1
      conditions {
        conditionable_type            = "AlertField"
        conditionable_id              = resource.rootly_alert_field.severity_field.id
        property_field_type           = "alert_field"
        property_field_condition_type = "is_one_of"
        property_field_values         = ["sev0", "sev1"]
      }
    }
  }

  rules {
    name = "Security Route"
    destinations {
      target_type = "EscalationPolicy"
      target_id   = rootly_escalation_policy.security_ep.id
    }
    condition_groups {
      position = 1
      conditions {
        conditionable_type            = "AlertField"
        conditionable_id              = resource.rootly_alert_field.service_field.id
        property_field_type           = "alert_field"
        property_field_condition_type = "contains"
        property_field_value          = "security"
      }
    }
  }

  rules {
    destinations {
      target_type = "EscalationPolicy"
      target_id   = rootly_escalation_policy.production_ep.id
    }
    fallback_rule = true
  }
}