package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAlertRoutingRule(t *testing.T) {
	t.Skip("Alert routing rule resource is deprecated - routing rules are now managed through alert routes")
	teamName := acctest.RandomWithPrefix("tf-team")
	alertsSourceName := acctest.RandomWithPrefix("tf-alerts-source")
	escalationPolicyName := acctest.RandomWithPrefix("tf-escalation-policy")
	alertRoutingRuleName := acctest.RandomWithPrefix("tf-alert-routing-rule")
	alertRoutingRuleNameUpdated := acctest.RandomWithPrefix("tf-alert-routing-rule-updated")
	resource.UnitTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertRoutingRuleConfig(teamName, alertsSourceName, escalationPolicyName, alertRoutingRuleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alert_routing_rule.test", "name", alertRoutingRuleName),
					resource.TestCheckResourceAttr("rootly_alert_routing_rule.test", "destination.0.target_type", "Group"),
				),
			},
			{
				Config: testAccResourceAlertRoutingRuleConfig(teamName, alertsSourceName, escalationPolicyName, alertRoutingRuleNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alert_routing_rule.test", "name", alertRoutingRuleNameUpdated),
					resource.TestCheckResourceAttr("rootly_alert_routing_rule.test", "destination.0.target_type", "Group"),
				),
			},
		},
	})
}

func testAccResourceAlertRoutingRuleConfig(teamName, alertsSourceName, escalationPolicyName, alertRoutingRuleName string) string {
	return fmt.Sprintf(`
resource "rootly_team" "test" {
	name = "%s"
}

resource "rootly_alerts_source" "test" {
  name = "%s"
  source_type = "generic_webhook"
  owner_group_ids = [rootly_team.test.id]
}

resource "rootly_escalation_policy" "test" {
  name      = "%s"
  group_ids = [rootly_team.test.id]
}

resource "rootly_alert_routing_rule" "test" {
  depends_on       = [rootly_escalation_policy.test]
  name             = "%s"
  alerts_source_id = rootly_alerts_source.test.id
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
`, teamName, alertsSourceName, escalationPolicyName, alertRoutingRuleName)
}
