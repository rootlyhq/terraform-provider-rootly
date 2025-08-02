package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAlertRoutingRule(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertRoutingRuleCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alert_routing_rule.test", "name", "Terraform"),
					resource.TestCheckResourceAttr("rootly_alert_routing_rule.test", "destination.0.target_type", "Group"),
				),
			},
			{
				Config: testAccResourceAlertRoutingRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alert_routing_rule.test", "name", "Terraform (updated)"),
					resource.TestCheckResourceAttr("rootly_alert_routing_rule.test", "destination.0.target_type", "Group"),
				),
			},
		},
	})
}

const testAccResourceAlertRoutingRuleCreate = `
data "rootly_alerts_source" "test" {
  source_type = "generic_webhook"
}

resource "rootly_team" "test" {
	name = "Test Team"
}

resource "rootly_escalation_policy" "test" {
  name      = "Test Terraform"
  group_ids = [rootly_team.test.id]
}

resource "rootly_alert_routing_rule" "test" {
  depends_on       = [rootly_escalation_policy.test]
  name             = "Terraform"
  alerts_source_id = data.rootly_alerts_source.test.id
  destination {
    target_id   = rootly_team.test.id
    target_type = "Group"
  }
  condition_type = "all"
  conditions {
    property_field_condition_type = "contains"
    property_field_name = "environment"
    property_field_type = "payload"
    property_field_value = "production"
  }
}
`

const testAccResourceAlertRoutingRuleUpdate = `
data "rootly_alerts_source" "test" {
  source_type = "generic_webhook"
}

resource "rootly_team" "test" {
	name = "Test Team"
}

resource "rootly_escalation_policy" "test" {
  name      = "Test Terraform"
  group_ids = [rootly_team.test.id]
}

resource "rootly_alert_routing_rule" "test" {
  depends_on       = [rootly_escalation_policy.test]
  name             = "Terraform (updated)"
  alerts_source_id = data.rootly_alerts_source.test.id
  destination {
    target_id   = rootly_team.test.id
    target_type = "Group"
  }
  condition_type = "all"
  conditions {
    property_field_condition_type = "contains"
    property_field_name = "environment"
    property_field_type = "payload"
    property_field_value = "production"
  }
}
`
