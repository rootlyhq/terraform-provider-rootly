resource "rootly_alerts_source" "example" {
  name        = "Generic webhook source"
  source_type = "generic_webhook"

  alert_template_attributes {
    title        = "Server exploded"
    description  = "Datacenter is burning down."
    external_url = "https://rootly.com"
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
    enabled = true
    condition_type = "and"
    identifier_json_path = "$.email.subject"
    identifier_value_regex = "ID:\\s*(\\w+)"
    conditions_attributes {
      field = "$.email.body"
      operator = "contains"
      value = "RESOLVED"
    }
    conditions_attributes {
      field = "$.email.body"
      operator = "does_not_contain"
      value = "ERROR"
    }
  }
}
