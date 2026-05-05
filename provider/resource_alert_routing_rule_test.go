package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAlertRoutingRule(t *testing.T) {
	t.Skip("Alert routing rule resource is deprecated - routing rules are now managed through alert routes")
	rName := acctest.RandomWithPrefix("tf-alert-rr")

	resource.UnitTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertRoutingRuleConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alert_routing_rule.test", "name", rName),
					resource.TestCheckResourceAttr("rootly_alert_routing_rule.test", "destination.0.target_type", "Group"),
				),
			},
			{
				Config: testAccResourceAlertRoutingRuleConfig(rName + "-updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alert_routing_rule.test", "name", rName+"-updated"),
					resource.TestCheckResourceAttr("rootly_alert_routing_rule.test", "destination.0.target_type", "Group"),
				),
			},
		},
	})
}

func testAccResourceAlertRoutingRuleConfig(rName string) string {
	return fmt.Sprintf(`
resource "rootly_team" "test" {
	name = "%s-team"
}

resource "rootly_alerts_source" "test" {
  name = "%s-source"
  source_type = "generic_webhook"
  owner_group_ids = [rootly_team.test.id]
}

resource "rootly_escalation_policy" "test" {
  name      = "%s-ep"
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
`, rName, rName, rName, rName)
}
