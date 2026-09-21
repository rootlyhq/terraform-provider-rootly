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
				Config: testAccResourceAlertGroupConfig(rName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alert_group.tf1", "owner_group_ids.#", "1"),
					resource.TestCheckResourceAttrPair("rootly_alert_group.tf1", "owner_group_ids.0", "rootly_team.tf", "id"),
				),
			},
			{
				Config: testAccResourceAlertGroupConfig(rName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alert_group.tf1", "owner_group_ids.#", "0"),
				),
			},
			{
				Config: testAccResourceAlertGroupConfig(rName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_alert_group.tf1", "owner_group_ids.#", "1"),
					resource.TestCheckResourceAttrPair("rootly_alert_group.tf1", "owner_group_ids.0", "rootly_team.tf", "id"),
				),
			},
		},
	})
}

func testAccResourceAlertGroupConfig(rName string, withOwnerTeam bool) string {
	ownerGroupIds := "[]"
	if withOwnerTeam {
		ownerGroupIds = "[rootly_team.tf.id]"
	}

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
	owner_group_ids = %s
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
`, rName, rName, rName, ownerGroupIds, rName)
}
