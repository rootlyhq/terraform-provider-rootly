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
				Config: testAccResourceAlertRoutingRule,
			},
		},
	})
}

const testAccResourceAlertRoutingRule = `
data "rootly_alerts_source" "terraform" {
  source_type = "generic_webhook"
}

resource "rootly_team" "terraform" {
	name = "test-routing-rule-tf"
}

resource "rootly_escalation_policy" "terraform" {
  name      = "Test Terraform"
  group_ids = [rootly_team.terraform.id]
}

resource "rootly_alert_routing_rule" "terraform" {
  depends_on       = [rootly_escalation_policy.terraform]
  name             = "Test Terraform"
  alerts_source_id = data.rootly_alerts_source.terraform.id
  destination {
    target_id   = rootly_team.terraform.id
    target_type = "Group"
  }
  condition_type = "all"
  conditions {
    property_field_condition_type = "is_one_of"
    property_field_name = "environment"
    property_field_type = "payload"
    property_field_value = "production"
  }
}
`
