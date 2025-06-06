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
				Config: testAccResourceAlertsSource,
			},
		},
	})
}

const testAccResourceAlertsSource = `
resource "rootly_team" "tf" {
	name = "tf"
	description = "tf"
}

resource "rootly_alert_urgency" "tf" {
	name = "tf"
	description = "tf"
	position = 1
}

resource "rootly_alerts_source" "tf" {
  name = "TF"

  source_type = "generic_webhook"
  owner_group_ids = [rootly_team.tf.id]

  alert_source_urgency_rules_attributes {
	alert_urgency_id = rootly_alert_urgency.tf.id
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
