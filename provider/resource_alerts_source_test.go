package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAlertsSource_UrgencyRulesDeletion(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-alert-source")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertsSourceWithUrgencyRule(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alerts_source.test", "name", rName),
					resource.TestCheckResourceAttr("rootly_alerts_source.test", "alert_source_urgency_rules_attributes.#", "1"),
					resource.TestCheckResourceAttr("rootly_alerts_source.test", "alert_source_urgency_rules_attributes.0.json_path", "$.severity"),
				),
			},
			{
				Config: testAccResourceAlertsSourceWithoutUrgencyRule(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alerts_source.test", "name", rName),
					resource.TestCheckResourceAttr("rootly_alerts_source.test", "alert_source_urgency_rules_attributes.#", "0"),
				),
			},
		},
	})
}

func testAccResourceAlertsSourceWithUrgencyRule(name string) string {
	return fmt.Sprintf(`
resource "rootly_alert_urgency" "test" {
	name        = "%[1]s-urgency"
	description = "Test urgency"
}

resource "rootly_alerts_source" "test" {
	name        = "%[1]s"
	source_type = "generic_webhook"

	alert_source_urgency_rules_attributes {
		json_path        = "$.severity"
		operator         = "is"
		value            = "critical"
		kind             = "payload"
		alert_urgency_id = rootly_alert_urgency.test.id
	}
}
`, name)
}

func testAccResourceAlertsSourceWithoutUrgencyRule(name string) string {
	return fmt.Sprintf(`
resource "rootly_alert_urgency" "test" {
	name        = "%[1]s-urgency"
	description = "Test urgency"
}

resource "rootly_alerts_source" "test" {
	name        = "%[1]s"
	source_type = "generic_webhook"
}
`, name)
}
