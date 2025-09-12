package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAlertRoute(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertRouteCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alert_route.test", "name", "Test Alert Route"),
					resource.TestCheckResourceAttr("rootly_alert_route.test", "enabled", "true"),
					resource.TestCheckResourceAttr("rootly_alert_route.test", "rules.0.name", "Production Alerts"),
					resource.TestCheckResourceAttr("rootly_alert_route.test", "rules.0.destinations.0.target_type", "Group"),
					resource.TestCheckResourceAttr("rootly_alert_route.test", "rules.0.condition_groups.0.conditions.0.property_field_condition_type", "contains"),
					resource.TestCheckResourceAttr("rootly_alert_route.test", "rules.0.condition_groups.0.conditions.0.property_field_name", "environment"),
					resource.TestCheckResourceAttr("rootly_alert_route.test", "rules.0.condition_groups.0.conditions.0.property_field_type", "payload"),
					resource.TestCheckResourceAttr("rootly_alert_route.test", "rules.0.condition_groups.0.conditions.0.property_field_value", "production"),
				),
			},
			{
				Config: testAccResourceAlertRouteUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alert_route.test", "name", "Updated Alert Route"),
					resource.TestCheckResourceAttr("rootly_alert_route.test", "enabled", "true"),
					resource.TestCheckResourceAttr("rootly_alert_route.test", "rules.0.name", "Critical Production Alerts"),
					resource.TestCheckResourceAttr("rootly_alert_route.test", "rules.0.destinations.0.target_type", "Group"),
					resource.TestCheckResourceAttr("rootly_alert_route.test", "rules.0.condition_groups.0.conditions.0.property_field_condition_type", "is_one_of"),
					resource.TestCheckResourceAttr("rootly_alert_route.test", "rules.0.condition_groups.0.conditions.0.property_field_name", "severity"),
					resource.TestCheckResourceAttr("rootly_alert_route.test", "rules.0.condition_groups.0.conditions.0.property_field_type", "attribute"),
				),
			},
		},
	})
}

func TestAccResourceAlertRouteWithComplexRules(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertRouteComplexRules,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alert_route.complex", "name", "Complex Alert Route"),
					resource.TestCheckResourceAttr("rootly_alert_route.complex", "rules.0.name", "Multi-Condition Rule"),
					resource.TestCheckResourceAttr("rootly_alert_route.complex", "rules.0.condition_groups.#", "2"),
					resource.TestCheckResourceAttr("rootly_alert_route.complex", "rules.0.condition_groups.0.conditions.#", "1"),
					resource.TestCheckResourceAttr("rootly_alert_route.complex", "rules.0.condition_groups.1.conditions.#", "1"),
					resource.TestCheckResourceAttr("rootly_alert_route.complex", "rules.1.name", "Fallback Rule"),
					resource.TestCheckResourceAttr("rootly_alert_route.complex", "rules.1.fallback_rule", "true"),
				),
			},
		},
	})
}

func TestAccResourceAlertRouteWithAlertField(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertRouteWithAlertField,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alert_route.with_field", "name", "Alert Route with Field"),
					resource.TestCheckResourceAttr("rootly_alert_route.with_field", "rules.0.condition_groups.0.conditions.0.property_field_type", "alert_field"),
					resource.TestCheckResourceAttrSet("rootly_alert_route.with_field", "rules.0.condition_groups.0.conditions.0.conditionable_id"),
				),
			},
		},
	})
}

const testAccResourceAlertRouteCreate = `
resource "rootly_team" "test" {
  name = "Test Team"
  description = "Test team for alert routing"
}

resource "rootly_alerts_source" "test" {
  name = "Test Alerts Source"
  source_type = "generic_webhook"
  owner_group_ids = [rootly_team.test.id]
}

resource "rootly_alert_route" "test" {
  name = "Test Alert Route"
  enabled = true
  alerts_source_ids = [rootly_alerts_source.test.id]
  owning_team_ids = [rootly_team.test.id]
  
  rules {
    name = "Production Alerts"
    position = 1
    fallback_rule = false
    
    destinations {
      target_type = "Group"
      target_id = rootly_team.test.id
    }
    
    condition_groups {
      position = 1
      conditions {
        property_field_condition_type = "contains"
        property_field_name = "environment"
        property_field_type = "payload"
        property_field_value = "production"
      }
    }
  }
}
`

const testAccResourceAlertRouteUpdate = `
resource "rootly_team" "test" {
  name = "Test Team"
  description = "Test team for alert routing"
}

resource "rootly_alerts_source" "test" {
  name = "Test Alerts Source"
  source_type = "generic_webhook"
  owner_group_ids = [rootly_team.test.id]
}

resource "rootly_alert_route" "test" {
  name = "Updated Alert Route"
  enabled = true
  alerts_source_ids = [rootly_alerts_source.test.id]
  owning_team_ids = [rootly_team.test.id]
  
  rules {
    name = "Critical Production Alerts"
    position = 1
    fallback_rule = false
    
    destinations {
      target_type = "Group"
      target_id = rootly_team.test.id
    }
    
    condition_groups {
      position = 1
      conditions {
        property_field_condition_type = "is_one_of"
        property_field_name = "severity"
        property_field_type = "attribute"
        property_field_values = ["critical", "high"]
      }
    }
  }
}
`

const testAccResourceAlertRouteComplexRules = `
resource "rootly_team" "test_primary" {
  name = "Primary Team"
  description = "Primary team for alerts"
}

resource "rootly_team" "test_fallback" {
  name = "Fallback Team"
  description = "Fallback team for alerts"
}

resource "rootly_alerts_source" "test" {
  name = "Test Alerts Source"
  source_type = "generic_webhook"
  owner_group_ids = [rootly_team.test_primary.id]
}

resource "rootly_alert_route" "complex" {
  name = "Complex Alert Route"
  enabled = true
  alerts_source_ids = [rootly_alerts_source.test.id]
  owning_team_ids = [rootly_team.test_primary.id]
  
  rules {
    name = "Multi-Condition Rule"
    position = 1
    fallback_rule = false
    
    destinations {
      target_type = "Group"
      target_id = rootly_team.test_primary.id
    }
    
    condition_groups {
      position = 1
      conditions {
        property_field_condition_type = "contains"
        property_field_name = "environment"
        property_field_type = "payload"
        property_field_value = "production"
      }
    }
    
    condition_groups {
      position = 2
      conditions {
        property_field_condition_type = "is_one_of"
        property_field_name = "severity"
        property_field_type = "attribute"
        property_field_values = ["critical", "high"]
      }
    }
  }
  
  rules {
    name = "Fallback Rule"
    position = 2
    fallback_rule = true
    
    destinations {
      target_type = "Group"
      target_id = rootly_team.test_fallback.id
    }
    
    condition_groups {
      position = 1
      conditions {
        property_field_condition_type = "is_empty"
        property_field_name = "assignee"
        property_field_type = "attribute"
      }
    }
  }
}
`

const testAccResourceAlertRouteWithAlertField = `
resource "rootly_team" "test" {
  name = "Test Team"
  description = "Test team for alert routing"
}

resource "rootly_alerts_source" "test" {
  name = "Test Alerts Source"
  source_type = "generic_webhook"
  owner_group_ids = [rootly_team.test.id]
}

resource "rootly_alert_field" "test" {
  name = "Test Alert Field"
}

resource "rootly_alert_route" "with_field" {
  name = "Alert Route with Field"
  enabled = true
  alerts_source_ids = [rootly_alerts_source.test.id]
  owning_team_ids = [rootly_team.test.id]
  
  rules {
    name = "Field-based Rule"
    position = 1
    fallback_rule = false
    
    destinations {
      target_type = "Group"
      target_id = rootly_team.test.id
    }
    
    condition_groups {
      position = 1
      conditions {
        property_field_condition_type = "contains"
        property_field_name = "custom_field"
        property_field_type = "alert_field"
        property_field_value = "important"
        conditionable_type = "AlertField"
        conditionable_id = rootly_alert_field.test.id
      }
    }
  }
}
`