package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAlertGroup(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertGroup,
			},
		},
	})
}

const testAccResourceAlertGroup = `
resource "rootly_alert_urgency" "tf" {
	name = "tf-test"
	description = "tf"
}

resource "rootly_team" "tf" {
	name = "tf-test"
}

resource "rootly_alert_group" "tf1" {
	name = "tf-test1"
	description = "tf"
	targets {
		target_type = "Group"
		target_id = rootly_team.tf.id
	}
	conditions {
		property_field_type = "payload"
		property_field_name = "monitor_id"
		property_field_condition_type = "matches_existing_alert"
	}
}

resource "rootly_alert_group" "tf2" {
	name = "tf-test2"
	description = "tf"
	conditions {
		property_field_type = "attribute"
		property_field_name = "alert_urgency"
		property_field_condition_type = "matches_existing_alert"
		values {
			record_id = rootly_alert_urgency.tf.id
			record_type = "AlertUrgency"
		}
	}
}
`
