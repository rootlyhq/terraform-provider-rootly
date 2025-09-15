package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAlertsSource(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertsSourceCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alerts_source.test", "name", "Test Alerts Source"),
					resource.TestCheckResourceAttr("rootly_alerts_source.test", "source_type", "generic_webhook"),
					resource.TestCheckResourceAttrSet("rootly_alerts_source.test", "id"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccResourceAlertsSourceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alerts_source.test", "name", "Test Alerts Source"),
					resource.TestCheckResourceAttr("rootly_alerts_source.test", "source_type", "generic_webhook"),
					resource.TestCheckResourceAttrSet("rootly_alerts_source.test", "id"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

const testAccResourceAlertsSourceCreate = `
resource "rootly_team" "tf" {
	name = "tf"
	description = "tf"
}

resource "rootly_alert_urgency" "tf" {
	name = "tf"
	description = "tf"
	position = 1
}

resource "rootly_alerts_source" "test" {
  name = "Test Alerts Source"

  source_type = "generic_webhook"
  owner_group_ids = [rootly_team.tf.id]

  alert_source_urgency_rules_attributes {
	alert_urgency_id = rootly_alert_urgency.tf.id
	json_path = "test"
	operator = "is"
	value = "P1"
  }

  sourceable_attributes {
    auto_resolve  = false
    resolve_state = "$.status"
  }
}
`

const testAccResourceAlertsSourceUpdate = `
resource "rootly_team" "tf" {
	name = "tf"
	description = "tf"
}

resource "rootly_alert_urgency" "tf" {
	name = "tf"
	description = "tf"
	position = 1
}

resource "rootly_alerts_source" "test" {
  name = "Test Alerts Source"

  source_type = "generic_webhook"
  owner_group_ids = [rootly_team.tf.id]

  alert_source_urgency_rules_attributes {
	alert_urgency_id = rootly_alert_urgency.tf.id
	json_path = "test"
	operator = "is"
	value = "P1"
  }

  sourceable_attributes {
    auto_resolve  = false
    resolve_state = "$.status"
  }
}
`
