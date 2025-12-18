# Required: Use data sources for title, description, and external_url fields
data "rootly_alert_field" "title_field" {
  kind = "title"
}

data "rootly_alert_field" "description_field" {
  kind = "description"
}

data "rootly_alert_field" "source_link_field" {
  kind = "external_url"
}

resource "rootly_alerts_source" "example" {
  name        = "Generic webhook source"
  source_type = "generic_webhook"

  # Required alert source fields: title, description, and external_url
  # template_body can be omitted if you want to use the default value for the alert source field
  alert_source_fields_attributes {
    alert_field_id = data.rootly_alert_field.title_field.id
    template_body  = "Server exploded"
  }

  alert_source_fields_attributes {
    alert_field_id = data.rootly_alert_field.description_field.id
    template_body  = "Datacenter is burning down."
  }

  alert_source_fields_attributes {
    alert_field_id = data.rootly_alert_field.source_link_field.id
    template_body  = "https://rootly.com"
  }

  sourceable_attributes {
    auto_resolve  = true
    resolve_state = "$.status"

    field_mappings_attributes {
      field     = "state"
      json_path = "$.my_group_attribute"
    }

    field_mappings_attributes {
      field     = "external_id"
      json_path = "$.my_id_attribute"
    }
  }

  resolution_rule_attributes {
    enabled                = true
    condition_type         = "and"
    identifier_json_path   = "$.email.subject"
    identifier_value_regex = "ID:\\s*(\\w+)"
    conditions_attributes {
      field    = "$.email.body"
      operator = "contains"
      value    = "RESOLVED"
    }
    conditions_attributes {
      field    = "$.email.body"
      operator = "does_not_contain"
      value    = "ERROR"
    }
  }
}
