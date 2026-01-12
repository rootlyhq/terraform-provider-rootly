resource "rootly_alert_route" "name_22" {
  name               = "name 22"
  enabled            = false
  alerts_source_ids  = ["ab2af6d8-5e35-4030-aeb7-39916aa268eb"]

}

resource "rootly_alert_route" "name" {
  name               = "name"
  enabled            = false
  alerts_source_ids  = ["ab2af6d8-5e35-4030-aeb7-39916aa268eb"]

}

resource "rootly_alert_route" "sailpoint_repro" {
  name               = "sailpoint repro"
  enabled            = false
  alerts_source_ids  = ["5cfc83f2-52d9-4cf1-8a91-849cbb31e942"]

  rules {
    name          = "1"
    position      = 1
    fallback_rule = false

    destinations {
      target_type = "EscalationPolicy"
      target_id   = "edc76590-b8f8-4476-962a-dae3e68c6be6"
    }


    condition_groups {
      position = 1
      conditions {
        property_field_type            = "alert_field"
        property_field_name            = "summary"
        property_field_condition_type  = "is_one_of"
        property_field_values          = ["Test"]
        conditionable_id               = "963ac299-8d35-4414-8bfa-aae34cd8d0e7"
        conditionable_type             = "AlertField"
      }
    }
  }
  rules {
    name          = "2"
    position      = 2
    fallback_rule = false

    destinations {
      target_type = "EscalationPolicy"
      target_id   = "33463e1b-3b11-4aeb-9b87-73e1bdefdf5d"
    }


    condition_groups {
      position = 1
      conditions {
        property_field_type            = "alert_field"
        property_field_name            = "summary"
        property_field_condition_type  = "starts_with"
        property_field_value           = "Test"
        conditionable_id               = "963ac299-8d35-4414-8bfa-aae34cd8d0e7"
        conditionable_type             = "AlertField"
      }
    }
  }
  rules {
    name          = "3"
    position      = 3
    fallback_rule = false

    destinations {
      target_type = "EscalationPolicy"
      target_id   = "db15fd8b-1775-438d-954b-fbfa71045ecc"
    }


    condition_groups {
      position = 1
      conditions {
        property_field_type            = "alert_field"
        property_field_name            = "external_url"
        property_field_condition_type  = "is_one_of"
        property_field_values          = ["google.com"]
        conditionable_id               = "b01eb0e8-6b10-4fb6-a2ef-66d322b1e682"
        conditionable_type             = "AlertField"
      }
    }
  }
  rules {
    name          = "4"
    position      = 4
    fallback_rule = false

    destinations {
      target_type = "EscalationPolicy"
      target_id   = "eb5ebb2c-ce23-42e5-9590-1e085ba23995"
    }


    condition_groups {
      position = 1
      conditions {
        property_field_type            = "alert_field"
        property_field_name            = "summary"
        property_field_condition_type  = "contains"
        property_field_value           = "TEST"
        conditionable_id               = "963ac299-8d35-4414-8bfa-aae34cd8d0e7"
        conditionable_type             = "AlertField"
      }
    }
  }
  rules {
    name          = "5"
    position      = 5
    fallback_rule = false

    destinations {
      target_type = "EscalationPolicy"
      target_id   = "8fee74d2-9a5c-4c85-905c-a086fe2817ec"
    }


    condition_groups {
      position = 1
      conditions {
        property_field_type            = "alert_field"
        property_field_name            = "summary"
        property_field_condition_type  = "contains"
        property_field_value           = "test"
        conditionable_id               = "963ac299-8d35-4414-8bfa-aae34cd8d0e7"
        conditionable_type             = "AlertField"
      }
    }
  }
  rules {
    name          = "Fallback Rule for sailpoint repro"
    position      = 6
    fallback_rule = true

    destinations {
      target_type = "EscalationPolicy"
      target_id   = "8d26c9f4-2596-437b-9358-78b98175493d"
    }

  }
}

resource "rootly_alert_route" "delete_source" {
  name               = "delete source"
  enabled            = false
  alerts_source_ids  = ["587d53ad-5e72-4836-8443-89e628a7aa0a"]

  rules {
    name          = "delete source"
    position      = 1
    fallback_rule = true

    destinations {
      target_type = "EscalationPolicy"
      target_id   = "edc76590-b8f8-4476-962a-dae3e68c6be6"
    }

  }
}

resource "rootly_alert_route" "shiny_new_routes_" {
  name               = "shiny new routes!"
  enabled            = false
  alerts_source_ids  = ["ee131642-9614-4fc3-9da3-33f2d133ef97", "c7ed6c17-871f-4601-ae67-56dd54351495", "2ee06f1c-0490-44bb-aaef-4bbdb0646f09", "fe32ddae-e645-475c-a0ca-eb7feabab2c1", "1ef96f68-237c-4333-8930-22c9ed8d1680", "eccf35d9-9988-445e-b654-5a5d860998b9", "238431fe-f0d7-418b-b719-6a95f353eb49"]

  rules {
    name          = "Rule #1"
    position      = 1
    fallback_rule = false

    destinations {
      target_type = "EscalationPolicy"
      target_id   = "4c3c3c2d-b080-45e6-8806-ed0240b87f7a"
    }


    condition_groups {
      position = 1
      conditions {
        property_field_type            = "attribute"
        property_field_name            = "alert_urgency"
        property_field_condition_type  = "is_one_of"
      }
    }
  }
  rules {
    name          = "added after fix"
    position      = 2
    fallback_rule = false

    destinations {
      target_type = "EscalationPolicy"
      target_id   = "1c1d25f3-a5ff-4055-957c-ce9902b2292b"
    }


    condition_groups {
      position = 1
      conditions {
        property_field_type            = "alert_field"
        property_field_name            = "🔥"
        property_field_condition_type  = "is_not_one_of"
        property_field_values          = ["hello"]
        conditionable_id               = "e5e9ea9c-25c2-47ad-8cf0-13cfa068e6fb"
        conditionable_type             = "AlertField"
      }
    }
  }
  rules {
    name          = "Fallback Rule for shiny new routes!"
    position      = 3
    fallback_rule = true

    destinations {
      target_type = "EscalationPolicy"
      target_id   = "eb5ebb2c-ce23-42e5-9590-1e085ba23995"
    }

  }
}

resource "rootly_alert_route" "new_route" {
  name               = "new route"
  enabled            = false
  alerts_source_ids  = ["ee131642-9614-4fc3-9da3-33f2d133ef97"]

  rules {
    name          = "rile"
    position      = 1
    fallback_rule = false

    destinations {
      target_type = "EscalationPolicy"
      target_id   = "1c1d25f3-a5ff-4055-957c-ce9902b2292b"
    }


    condition_groups {
      position = 1
      conditions {
        property_field_type            = "alert_field"
        property_field_name            = "Total Incidents"
        property_field_condition_type  = "contains"
        property_field_value           = "fana"
        conditionable_id               = "448cf823-179a-4622-8e3c-ccf7414f1419"
        conditionable_type             = "AlertField"
      }
    }
  }
  rules {
    name          = "new route"
    position      = 2
    fallback_rule = true

    destinations {
      target_type = "EscalationPolicy"
      target_id   = "eb5ebb2c-ce23-42e5-9590-1e085ba23995"
    }

  }
}

resource "rootly_alert_route" "uppercase_route" {
  name               = "UPPERCASE ROUTE"
  enabled            = false
  alerts_source_ids  = ["ee131642-9614-4fc3-9da3-33f2d133ef97"]

  rules {
    name          = "UPPERCASE ROUTE"
    position      = 2
    fallback_rule = true

    destinations {
      target_type = "EscalationPolicy"
      target_id   = "8fee74d2-9a5c-4c85-905c-a086fe2817ec"
    }

    destinations {
      target_type = "EscalationPolicy"
      target_id   = "db15fd8b-1775-438d-954b-fbfa71045ecc"
    }

  }
}

resource "rootly_alert_route" "new_relic_route" {
  name               = "new relic route"
  enabled            = false
  alerts_source_ids  = ["87e143b3-10af-434c-84b3-bdd0ed61abed"]

  rules {
    name          = "new relic route"
    position      = 1
    fallback_rule = true

    destinations {
      target_type = "EscalationPolicy"
      target_id   = "8fee74d2-9a5c-4c85-905c-a086fe2817ec"
    }

  }
}

resource "rootly_alert_route" "new_new_generic_route" {
  name               = "new new generic route"
  enabled            = false
  alerts_source_ids  = ["238431fe-f0d7-418b-b719-6a95f353eb49"]

  rules {
    name          = "Fallback Rule for new new generic route"
    position      = 1
    fallback_rule = true

    destinations {
      target_type = "EscalationPolicy"
      target_id   = "1c1d25f3-a5ff-4055-957c-ce9902b2292b"
    }

  }
}

resource "rootly_alert_route" "grafana_route" {
  name               = "grafana route"
  enabled            = true
  alerts_source_ids  = ["ee131642-9614-4fc3-9da3-33f2d133ef97"]

  rules {
    name          = "grafana route"
    position      = 1
    fallback_rule = true

    destinations {
      target_type = "EscalationPolicy"
      target_id   = "8fee74d2-9a5c-4c85-905c-a086fe2817ec"
    }

  }
}

resource "rootly_alert_route" "inc_450_repro_route" {
  name               = "inc 450 repro route"
  enabled            = false
  alerts_source_ids  = ["e204225c-7caf-467b-b676-6ab77c11a3a0"]

  rules {
    name          = "inc 450 repro route"
    position      = 1
    fallback_rule = true

    destinations {
      target_type = "EscalationPolicy"
      target_id   = "8fee74d2-9a5c-4c85-905c-a086fe2817ec"
    }

  }
}

resource "rootly_alert_route" "route_via_web" {
  name               = "route via web"
  enabled            = false
  alerts_source_ids  = ["20eab254-1c76-446a-afed-94aee8d8cbf8", "87e143b3-10af-434c-84b3-bdd0ed61abed", "ab2af6d8-5e35-4030-aeb7-39916aa268eb"]

  rules {
    name          = "Fallback Rule for route via web"
    position      = 1
    fallback_rule = true

    destinations {
      target_type = "EscalationPolicy"
      target_id   = "8fee74d2-9a5c-4c85-905c-a086fe2817ec"
    }

  }
}

resource "rootly_alert_route" "new_alert" {
  name               = "new alert"
  enabled            = false
  alerts_source_ids  = ["f28990df-63ef-4f09-bb2c-45231c807d12"]

}

resource "rootly_alert_route" "webhook_source" {
  name               = "webhook source"
  enabled            = true
  alerts_source_ids  = ["238431fe-f0d7-418b-b719-6a95f353eb49"]

  rules {
    name          = "webhook source"
    position      = 1
    fallback_rule = true

    destinations {
      target_type = "EscalationPolicy"
      target_id   = "4320c378-1a49-4801-aa3f-dd7380e92c1f"
    }

  }
}

resource "rootly_alert_route" "generic_wehbook_source" {
  name               = "generic wehbook source"
  enabled            = false
  alerts_source_ids  = ["238431fe-f0d7-418b-b719-6a95f353eb49"]
  owning_team_ids    = ["4e561797-bb24-4ca3-a536-b650d0a03bbf"]

  rules {
    name          = "generic wehbook source"
    position      = 1
    fallback_rule = true

    destinations {
      target_type = "EscalationPolicy"
      target_id   = "db15fd8b-1775-438d-954b-fbfa71045ecc"
    }

  }
}


# Import statements
# terraform import rootly_alert_route.name_22 'd379a671-42c0-4e89-868c-b6e456f406b8'
# terraform import rootly_alert_route.name 'e8957498-8e1c-498f-b809-82fd6d195812'
# terraform import rootly_alert_route.sailpoint_repro '0a39d513-5e03-4779-ab44-cabd95e90218'
# terraform import rootly_alert_route.delete_source 'f98d9d47-744d-4389-94a4-ffec8295600b'
# terraform import rootly_alert_route.shiny_new_routes_ 'da9239d7-0fc4-4da1-b9d8-8359fb10cc62'
# terraform import rootly_alert_route.new_route '94503850-476e-44d9-b5c8-7aaef49e775e'
# terraform import rootly_alert_route.uppercase_route '67395ca7-e98c-48bc-b11e-95b949a98ad7'
# terraform import rootly_alert_route.new_relic_route 'ffe8587a-5aac-4fe7-b42d-09d44cdb8bca'
# terraform import rootly_alert_route.new_new_generic_route 'ce81a55b-d1f4-4fcf-ae12-3683b5faf780'
# terraform import rootly_alert_route.grafana_route '9f76d75a-fc7d-41fa-9e4e-fc9f0a900f12'
# terraform import rootly_alert_route.inc_450_repro_route '2fdf4bd5-985e-4ccc-845c-b6c0aadabe1c'
# terraform import rootly_alert_route.route_via_web '64bb9757-ac40-47f6-9c82-17587cb4a43c'
# terraform import rootly_alert_route.new_alert 'd5a3bcf8-e13f-44de-9962-4c1bd3712ec7'
# terraform import rootly_alert_route.webhook_source 'df330656-41ce-4860-baec-96fa37863b7f'
# terraform import rootly_alert_route.generic_wehbook_source 'fdb6e953-c535-4dd7-a413-fac302836abc'

# NOTE: If you prefer import blocks over import statements, re-run with -import=block flag
#
# Instructions:
# 1. Run the terraform import commands above to import existing resources into state
# 2. Run 'terraform plan' to verify the resources are properly imported
#    WARNING: You should see NO changes in the plan - no creates, updates, deletes or imports
# 3. Remove the import statements above from this file once imports are complete
# 4. Run 'terraform state rm <resource_address>' for each deprecated alert_routing_rule resource
#    Example: terraform state rm rootly_alert_routing_rule.my_rule
# 5. Remove deprecated 'rootly_alert_routing_rule' resources from your Terraform configuration
