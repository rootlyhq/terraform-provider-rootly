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
				),
			},
			{
				Config: testAccResourceAlertRouteUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alert_route.test", "name", "Updated Alert Route"),
					resource.TestCheckResourceAttr("rootly_alert_route.test", "enabled", "false"),
				),
			},
		},
	})
}

const testAccResourceAlertRouteCreate = `
resource "rootly_team" "test" {
	name = "test-team"
	description = "Test team for alert route"
}

resource "rootly_alerts_source" "test" {
	name = "test-alerts-source"
	source_type = "generic_webhook"
	owner_group_ids = [rootly_team.test.id]
}

resource "rootly_alert_route" "test" {
	name = "Test Alert Route"
	enabled = true
	alerts_source_ids = [rootly_alerts_source.test.id]
	owner_group_ids = [rootly_team.test.id]
	alert_routing_rules_attributes {
		name = "Test Routing Rule"
		enabled = true
		condition_type = "all"
		alert_routing_rule_targets {
			target_type = "Group"
			target_id = rootly_team.test.id
		}
		alert_routing_rule_condition_groups {
			alert_routing_rule_conditions {
				property_field_type = "payload"
				property_field_name = "severity"
				property_field_condition_type = "contains"
				property_field_value = "high"
			}
		}
	}
}
`

const testAccResourceAlertRouteUpdate = `
resource "rootly_team" "test" {
	name = "test-team"
	description = "Test team for alert route"
}

resource "rootly_alerts_source" "test" {
	name = "test-alerts-source"
	source_type = "generic_webhook"
	owner_group_ids = [rootly_team.test.id]
}

resource "rootly_alert_route" "test" {
	name = "Updated Alert Route"
	enabled = false
	alerts_source_ids = [rootly_alerts_source.test.id]
	owner_group_ids = [rootly_team.test.id]
	alert_routing_rules_attributes {
		name = "Updated Routing Rule"
		enabled = true
		condition_type = "any"
		alert_routing_rule_targets {
			target_type = "Group"
			target_id = rootly_team.test.id
		}
		alert_routing_rule_condition_groups {
			alert_routing_rule_conditions {
				property_field_type = "payload"
				property_field_name = "severity"
				property_field_condition_type = "contains"
				property_field_value = "critical"
			}
		}
	}
}
`
