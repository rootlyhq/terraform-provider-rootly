package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAlertGroup(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-alert-group")

	resource.UnitTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertGroupConfig(rName),
			},
		},
	})
}

func testAccResourceAlertGroupConfig(rName string) string {
	return fmt.Sprintf(`
resource "rootly_alert_urgency" "tf" {
	name = "%s-urgency"
	description = "tf"
}

resource "rootly_team" "tf" {
	name = "%s-team"
}

resource "rootly_alert_group" "tf1" {
	name = "%s-1"
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
	name = "%s-2"
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
`, rName, rName, rName, rName)
}
