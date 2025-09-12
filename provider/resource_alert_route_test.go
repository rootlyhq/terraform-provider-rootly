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
					resource.TestCheckResourceAttrSet("rootly_alert_route.test", "id"),
				),
			},
			{
				Config: testAccResourceAlertRouteUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alert_route.test", "name", "Updated Alert Route"),
					resource.TestCheckResourceAttr("rootly_alert_route.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("rootly_alert_route.test", "id"),
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
					resource.TestCheckResourceAttr("rootly_alert_route.complex", "enabled", "true"),
					resource.TestCheckResourceAttrSet("rootly_alert_route.complex", "id"),
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
					resource.TestCheckResourceAttr("rootly_alert_route.with_field", "enabled", "true"),
					resource.TestCheckResourceAttrSet("rootly_alert_route.with_field", "id"),
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

resource "rootly_escalation_policy" "test" {
  name = "Test Escalation Policy"
  group_ids = [rootly_team.test.id]
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
      target_type = "EscalationPolicy"
      target_id = rootly_escalation_policy.test.id
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

resource "rootly_escalation_policy" "test" {
  name = "Test Escalation Policy"
  group_ids = [rootly_team.test.id]
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
      target_type = "EscalationPolicy"
      target_id = rootly_escalation_policy.test.id
    }
    
    condition_groups {
      position = 1
      conditions {
        property_field_condition_type = "contains"
        property_field_name = "summary"
        property_field_type = "attribute"
        property_field_value = "critical"
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

resource "rootly_escalation_policy" "test_primary" {
  name = "Primary Escalation Policy"
  group_ids = [rootly_team.test_primary.id]
}

resource "rootly_escalation_policy" "test_fallback" {
  name = "Fallback Escalation Policy"
  group_ids = [rootly_team.test_fallback.id]
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
      target_type = "EscalationPolicy"
      target_id = rootly_escalation_policy.test_primary.id
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
        property_field_condition_type = "contains"
        property_field_name = "summary"
        property_field_type = "attribute"
        property_field_value = "critical"
      }
    }
  }
  
  rules {
    name = "Fallback Rule"
    position = 2
    fallback_rule = true
    
    destinations {
      target_type = "EscalationPolicy"
      target_id = rootly_escalation_policy.test_fallback.id
    }
    
    condition_groups {
      position = 1
      conditions {
        property_field_condition_type = "contains"
        property_field_name = "description"
        property_field_type = "attribute"
        property_field_value = "fallback"
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

resource "rootly_escalation_policy" "test" {
  name = "Test Escalation Policy"
  group_ids = [rootly_team.test.id]
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
      target_type = "EscalationPolicy"
      target_id = rootly_escalation_policy.test.id
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