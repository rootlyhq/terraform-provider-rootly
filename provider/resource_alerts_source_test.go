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
					resource.TestCheckResourceAttrSet("rootly_alerts_source.with-template", "id"),
				),
			},
			{
				Config: testAccResourceAlertsSourceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("rootly_alerts_source.with-template", "id"),
				),
			},
		},
	})
}

const testAccResourceAlertsSourceCreate = `
resource "rootly_team" "test" {
	name = "Test Team Alerts Source"
	description = "tf"
}

resource "rootly_alert_urgency" "test" {
	name = "Test Urgency Alerts Source"
	description = "tf"
	position = 1
}

resource "rootly_alerts_source" with-template {
  name = "TF: with template"

  source_type = "generic_webhook"
  owner_group_ids = [rootly_team.test.id]

  alert_source_urgency_rules_attributes {
	alert_urgency_id = rootly_alert_urgency.test.id
	json_path = "test"
	operator = "is"
	value = "P1"
  }

  alert_template_attributes {
    title        = "Server exploded"
    description  = "Datacenter is burning down."
    external_url = "https://rootly.com"
  }

  sourceable_attributes {
    auto_resolve  = false
    resolve_state = "$.status"
  }
}
`

const testAccResourceAlertsSourceUpdate = `
resource "rootly_team" "test" {
	name = "Test Team Alerts Source"
	description = "tf"
}

resource "rootly_alert_urgency" "test" {
	name = "Test Urgency Alerts Source"
	description = "tf"
	position = 1
}

resource "rootly_alerts_source" with-template {
  name = "TF: with template"

  source_type = "generic_webhook"
  owner_group_ids = [rootly_team.test.id]

  alert_source_urgency_rules_attributes {
	alert_urgency_id = rootly_alert_urgency.test.id
	json_path = "test"
	operator = "is"
	value = "P1"
  }

  alert_template_attributes {
    title        = "Server exploded!"
    description  = "Datacenter is burning down."
    external_url = "https://rootly.com"
  }

  sourceable_attributes {
    auto_resolve  = false
    resolve_state = "$.status"
  }
}
`
